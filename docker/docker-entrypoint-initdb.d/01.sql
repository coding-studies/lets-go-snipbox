--
-- Make sure this role is not the same as the root PG role set in .env,
-- otherwise, the setup scripts will think the data directory is not empty and
-- the SQL code in this file will not be executed. It will complain that the
-- role 'dev' already exists and stop, and the setup below will not happen.
--
CREATE ROLE dev
WITH LOGIN PASSWORD 'dev'
CREATEDB REPLICATION
VALID UNTIL 'infinity';

----
-- When we log into a PostgreSQL server, we are required to “be” in some
-- database. Let's create a `dev` DB for the `dev` role so when we log in with
-- the `dev` role, it will land in the `dev` DB by default if a specific DB name
-- is not provided.
--
CREATE DATABASE dev WITH
    ENCODING='UTF8'
    OWNER=dev
    LC_CTYPE='en_US.UTF-8'
    LC_COLLATE='en_US.UTF-8'
    TEMPLATE=template0
    CONNECTION LIMIT=7;

CREATE DATABASE snipbox_dev WITH
    ENCODING='UTF8'
    OWNER=dev
    LC_CTYPE='en_US.UTF-8'
    LC_COLLATE='en_US.UTF-8'
    TEMPLATE=template0
    CONNECTION LIMIT=7;

\c snipbox_dev;

CREATE TABLE users (
    id SERIAL PRIMARY KEY
  , name VARCHAR(128) NOT NULL
);

----
-- From the docs:
--
--   CREATE TABLE will create a new, initially empty table in the current
--   database. The table will be owned by the user issuing the command.
--
-- https://www.postgresql.org/docs/17/sql-createtable.html
--
-- These SQL commands are being run with the user defined in the `.env` file.
-- Therefore, let's update the owner to the one we want.
--
ALTER TABLE users OWNER TO dev;

----
-- We must not provide the IDs explicitly. We must instead let PostgreSQL handle
-- it so future insertions will continue incrementing the IDs automatially. If
-- we insert the IDs ourselves, like 1, 2, 3, then future insertions with
-- automatic ID will try to start from 1, and it will fail as there are already
-- the IDs 1, 2 and 3 in the table.
--
INSERT INTO users (
  name
) VALUES
    ('Ahsoka Tano')
  , ('Leia Organa')
  , ('Aayla Secura');

CREATE TABLE IF NOT EXISTS snippets (
                                      id SERIAL PRIMARY KEY
  , title VARCHAR(128) NOT NULL
  , content TEXT NOT NULL
  , created TIMESTAMP NOT NULL
  , expires TIMESTAMP NOT NULL
  );

ALTER TABLE users OWNER TO dev;

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (
    title
  , content
  , created
  , expires
) VALUES (
    'An old silent pond'
  , 'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō'
  , CURRENT_TIMESTAMP
  , CURRENT_TIMESTAMP + '1 year'
), (
     'Over the wintry forest'
   , 'Over the wintry forest\n,winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki'
   , CURRENT_TIMESTAMP
   , CURRENT_TIMESTAMP + '1 year'
), (
     'First autumn morning'
   , 'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo'
   , CURRENT_TIMESTAMP
   , CURRENT_TIMESTAMP + '1 year'
);