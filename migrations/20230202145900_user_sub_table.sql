-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS subscriptions (
    id uuid default uuid_generate_v4() not null
        constraint subscriptions_pkey
            primary key,
    chat_id integer NOT NULL,
    city character varying(255) NOT NULL,
    district character varying(255) NOT NULL,
    user_name character varying(255) NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS subscriptions;
DROP EXTENSION IF EXISTS "uuid-ossp";
