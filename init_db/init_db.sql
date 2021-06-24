CREATE USER forum_root WITH password 'root';
-- CREATE EXTENSION citext;

DROP DATABASE IF EXISTS forum;
CREATE DATABASE forum
    WITH OWNER forum_root
    ENCODING 'utf8';
GRANT ALL PRIVILEGES ON DATABASE forum TO forum_root;
\connect forum;

CREATE EXTENSION citext;

DROP TABLE IF EXISTS users CASCADE;
CREATE UNLOGGED TABLE users (
    nickname CITEXT UNIQUE PRIMARY KEY,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);
CREATE UNIQUE INDEX ON users(email);

GRANT ALL PRIVILEGES ON TABLE users TO forum_root;


DROP TABLE IF EXISTS forums CASCADE;
CREATE UNLOGGED TABLE forums (
   title CITEXT NOT NULL,
   "user" CITEXT UNIQUE REFERENCES users(nickname),
   slug CITEXT UNIQUE NOT NULL PRIMARY KEY,
   posts BIGINT DEFAULT 0,
   threads INT DEFAULT 0
);
GRANT ALL PRIVILEGES ON TABLE forums TO forum_root;

DROP TABLE IF EXISTS threads CASCADE;
CREATE UNLOGGED TABLE threads (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author CITEXT NOT NULL REFERENCES users(nickname),
    forum CITEXT NOT NULL REFERENCES forums(slug),
    message TEXT NOT NULL,
    votes INT DEFAULT 0,
    slug CITEXT,
    created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX ON threads (slug) WHERE slug IS NOT NULL;
CREATE INDEX ON threads (forum, created);
CREATE INDEX ON threads (created);
GRANT ALL PRIVILEGES ON TABLE threads TO forum_root;

DROP TABLE IF EXISTS voices CASCADE;
CREATE UNLOGGED TABLE voices (
    nickname CITEXT REFERENCES users(nickname),
    voice INT,
    thread INT REFERENCES threads(id)
);
CREATE UNIQUE INDEX ON voices(nickname, thread);
GRANT ALL PRIVILEGES ON TABLE voices TO forum_root;

DROP TABLE IF EXISTS posts CASCADE;
CREATE UNLOGGED TABLE posts (
    author CITEXT NOT NULL REFERENCES users(nickname),
    created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum CITEXT NOT NULL REFERENCES forums(slug),
    id BIGSERIAL PRIMARY KEY,
    is_edited BOOLEAN DEFAULT false,
    message TEXT NOT NULL,
    parent BIGINT DEFAULT 0,
    thread INT REFERENCES threads(id),
    path BIGINT[] DEFAULT ARRAY []::INTEGER[]
);
CREATE INDEX ON posts((path[1]));
CREATE INDEX ON posts(id, (path[1]));
CREATE UNIQUE INDEX ON posts(id, thread);
CREATE UNIQUE INDEX ON posts(id, author);
CREATE INDEX ON posts(thread, path, id);
CREATE INDEX ON posts(thread, id);

GRANT ALL PRIVILEGES ON TABLE posts TO forum_root;

CREATE OR REPLACE FUNCTION insert_votes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE threads SET votes=(votes + NEW.voice) WHERE id=NEW.thread;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_votes() RETURNS TRIGGER AS
$update_users_forum$
begin
    IF OLD.voice <> NEW.voice THEN
        UPDATE threads SET votes=(votes + NEW.voice - OLD.voice) WHERE id=NEW.thread;
    END IF;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_vote
    BEFORE INSERT
    ON voices
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes();

CREATE TRIGGER trigger_update_vote
    BEFORE UPDATE
    ON voices
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();

CREATE FUNCTION insert_posts() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums SET
    posts = posts + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_posts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE insert_posts();

CREATE FUNCTION insert_threads() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums SET
        threads = threads + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE insert_threads();


CREATE FUNCTION update_path() RETURNS TRIGGER AS $$
DECLARE
    parent_thread INTEGER;
BEGIN
    IF NEW.parent = 0 THEN
        NEW.path = ARRAY [NEW.id];
    ELSE
        SELECT thread
        INTO parent_thread
        FROM posts
        WHERE id = new.parent;

        IF parent_thread ISNULL THEN
            RAISE EXCEPTION 'Parent post not found %', NEW.parent;
        ELSIF parent_thread <> NEW.thread THEN
            RAISE EXCEPTION 'Thread not found %', NEW.thread;
        END IF;

        SELECT path
        INTO NEW.path
        FROM posts
        WHERE id = NEW.parent;
        NEW.path = array_append(NEW.path, NEW.id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();


GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO forum_root;


