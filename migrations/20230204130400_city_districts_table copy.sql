-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "fuzzystrmatch";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE TABLE IF NOT EXISTS districts (
    id uuid default uuid_generate_v4() not null
        constraint districts_pkey
            primary key,
    district character varying(255) NOT NULL,
    city character varying(255) NOT NULL
   
);

INSERT INTO districts (district, city)
VALUES ('Babataşı','Fethiye'),
('Karaçulha','Fethiye'),
('Patlangıç','Fethiye'),
('Menteşeoğlu','Fethiye'),
('Akarca','Fethiye'),
('Göcek','Fethiye'),
('Çiftlik','Fethiye'), 
('Cumhuriyet','Fethiye'),
('Yeşilüzümlü','Fethiye'),
('Karagözler','Fethiye'),
('İnlice','Fethiye'),
('Kayaköy','Fethiye'),
('Gökben','Fethiye'),
('Karaağaç','Fethiye'),
('Dalyan','Ortaca'),
('Cumhuriyet','Ortaca'),
('Terzialiler','Ortaca'),
('Ataturk','Ortaca'),
('Bahcelievler','Ortaca'),
('Beskopru','Ortaca'),
('Karaburun','Ortaca'),
('Eksiliyurt','Ortaca'),
('Cayli','Ortaca')
('Foça', 'Fethiye');


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS districts;
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "fuzzystrmatch";
DROP EXTENSION IF EXISTS "pg_trgm";
