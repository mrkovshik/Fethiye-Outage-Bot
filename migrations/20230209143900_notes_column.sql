-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE IF EXISTS outages ADD COLUMN notes character varying(255);
ALTER TABLE IF EXISTS outages ADD COLUMN date_added timestamp without time zone;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.


