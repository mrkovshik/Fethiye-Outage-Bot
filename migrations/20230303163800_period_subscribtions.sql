-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE subscriptions 
ADD column IF NOT EXISTS period integer NOT NULL,
DROP COLUMN IF exists user_name;
ALTER TABLE subscriptions RENAME TO subscribtions;



-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE subscribtions 
DROP column IF EXISTS period,
ADD COLUMN IF NOT exists user_name TEXT;
ALTER TABLE subscribtions RENAME TO subscriptions;
