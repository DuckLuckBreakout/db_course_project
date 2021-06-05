CREATE USER forum_root WITH password 'root';
CREATE EXTENSION citext;

DROP DATABASE IF EXISTS forum;
CREATE DATABASE forum
    WITH OWNER forum_root
    ENCODING 'utf8';
GRANT ALL PRIVILEGES ON DATABASE forum TO forum_root;
\connect forum;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    nickname CITEXT UNIQUE PRIMARY KEY,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);
GRANT ALL PRIVILEGES ON TABLE users TO forum_root;


DROP TABLE IF EXISTS forums CASCADE;
CREATE TABLE forums (
   title CITEXT UNIQUE NOT NULL,
   "user" CITEXT UNIQUE REFERENCES users(nickname),
   slug CITEXT UNIQUE NOT NULL PRIMARY KEY,
   posts BIGINT DEFAULT 0,
   threads INT DEFAULT 0
);
GRANT ALL PRIVILEGES ON TABLE forums TO forum_root;

DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE threads (
    id INT UNIQUE PRIMARY KEY,
    title TEXT NOT NULL,
    author CITEXT NOT NULL REFERENCES users(nickname),
    forum CITEXT UNIQUE NOT NULL REFERENCES forums(slug),
    message TEXT NOT NULL,
    votes INT,
    slug CITEXT,
    created DATE
);
GRANT ALL PRIVILEGES ON TABLE threads TO forum_root;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO forum_root;

