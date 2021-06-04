CREATE USER forum_root WITH password 'root';
CREATE EXTENSION citext;

DROP DATABASE IF EXISTS forum;
CREATE DATABASE forum
    WITH OWNER forum_root
    ENCODING 'utf8';
GRANT ALL PRIVILEGES ON database forum TO forum_root;
\connect forum;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    nickname CITEXT UNIQUE,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);
GRANT ALL PRIVILEGES ON TABLE users TO forum_root;


GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO forum_root;

