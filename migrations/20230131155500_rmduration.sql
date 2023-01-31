-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE outages
DROP COLUMN IF EXISTS duration;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
