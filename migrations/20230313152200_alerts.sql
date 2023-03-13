-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS alerts (
    id uuid default uuid_generate_v4() not null
        constraint alerts_pkey
            primary key,
    is_sent bool NOT NULL,
    cancelled bool NOT NULL,
    alert_time timestamp without time zone NOT NULL,
    outage_id uuid REFERENCES outages(id),
    subscribtion_id uuid REFERENCES subscribtions(id)
);

ALTER TABLE subscribtions 
ADD column IF NOT EXISTS cancelled bool NOT NULL,
ADD column IF NOT EXISTS city_normalized TEXT,
ADD COLUMN IF NOT exists district_normalized TEXT;
ALTER TABLE outages 
ADD column IF NOT EXISTS alerted bool NOT NULL  default false;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS alerts;
DROP EXTENSION IF EXISTS "uuid-ossp";
ALTER TABLE subscriptions 
DROP column IF  EXISTS cancelled bool NOT NULL,
DROP column IF  EXISTS city_normalized TEXT,
DROP COLUMN IF  exists district_normalized TEXT;
ALTER TABLE outages 
DROP column IF EXISTS alerted bool NOT NULL;