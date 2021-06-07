CREATE USER forum_root WITH password 'root';
CREATE EXTENSION citext;

DROP DATABASE IF EXISTS forum;
CREATE DATABASE forum
    WITH OWNER forum_root
    ENCODING 'utf8';
GRANT ALL PRIVILEGES ON DATABASE forum TO forum_root;
\connect forum;

DROP TABLE IF EXISTS users CASCADE;
CREATE UNLOGGED TABLE users (
    nickname CITEXT UNIQUE PRIMARY KEY,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);
GRANT ALL PRIVILEGES ON TABLE users TO forum_root;


DROP TABLE IF EXISTS forums CASCADE;
CREATE UNLOGGED TABLE forums (
   title CITEXT UNIQUE NOT NULL,
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
GRANT ALL PRIVILEGES ON TABLE threads TO forum_root;

DROP TABLE IF EXISTS voices CASCADE;
CREATE UNLOGGED TABLE voices (
    nickname CITEXT REFERENCES users(nickname),
    voice INT,
    thread INT REFERENCES threads(id)
);
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

CREATE TRIGGER add_vote
    BEFORE INSERT
    ON voices
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes();

CREATE TRIGGER edit_vote
    BEFORE UPDATE
    ON voices
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();


GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO forum_root;

