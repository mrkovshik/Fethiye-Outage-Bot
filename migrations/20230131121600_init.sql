-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS outages (
    id uuid default uuid_generate_v4() not null
        constraint outages_pkey
            primary key,
    resource character varying(255) NOT NULL,
    city character varying(255) NOT NULL,
    district character varying(255) NOT NULL,
    start_date timestamp without time zone NOT NULL,
    duration integer NOT NULL,
    end_date timestamp without time zone NOT NULL,
    source_url character varying(255) NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS outages;
DROP EXTENSION IF EXISTS "uuid-ossp";
