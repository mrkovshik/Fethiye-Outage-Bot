-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "fuzzystrmatch";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS outages (
    id uuid default uuid_generate_v4() not null
        constraint outages_pkey
            primary key,
    resource character varying(255) NOT NULL,
    city character varying(255) NOT NULL,
    district character varying(255) NOT NULL,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    source_url character varying(255) NOT NULL,
    date_added timestamp without time zone NOT NULL,
    notes text NOT NULL
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id uuid default uuid_generate_v4() not null
        constraint subscriptions_pkey
            primary key,
    chat_id integer NOT NULL,
    city character varying(255) NOT NULL,
    district character varying(255) NOT NULL,
    user_name character varying(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS districts (
    id uuid default uuid_generate_v4() not null
        constraint districts_pkey
            primary key,
    district character varying(255) NOT NULL,
    city character varying(255) NOT NULL
);

INSERT INTO outages (resource, city, district, start_date, end_date, source_url, notes, date_added)
VALUES ('water','Limpopo','Ugadagada','2023-02-03 18:37:56', '2123-02-03 18:37:56', 'test entry','test entry', '2123-02-03 18:37:56'),
('water','Limpopo','Butumaputu','2023-02-03 18:37:56', '2123-02-03 18:37:56', 'test entry','test entry', '2123-02-03 18:37:56'),
('water','Limpopo','Ugadagada','2023-02-03 18:37:56', '2023-02-03 20:37:56', 'test entry','test entry', '2123-02-03 18:37:56')
;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS districts;
DROP TABLE IF EXISTS outages;
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "fuzzystrmatch";
DROP EXTENSION IF EXISTS "pg_trgm"
