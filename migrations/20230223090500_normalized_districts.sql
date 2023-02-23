-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;


ALTER TABLE districts 
ADD column IF NOT EXISTS city_normalized TEXT,
ADD COLUMN IF NOT exists district_normalized TEXT;

UPDATE districts 
SET city_normalized = lower (normalize(regexp_replace(unaccent(city), '[^a-zA-Z0-9\s]+', ' ', 'g'), NFD)),
district_normalized = lower (normalize(regexp_replace(unaccent(district), '[^a-zA-Z0-9\s]+', ' ', 'g'),NFD));

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP EXTENSION IF EXISTS "pg_trgm",
DROP EXTENSION IF EXISTS "unaccent";

DROP column IF EXISTS city_normalized,
DROP column IF EXISTS district_normalized,
