CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL PRIMARY KEY,
    login text NOT NULL UNIQUE,
    password bytea NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamp default (now() at time zone 'utc')
);