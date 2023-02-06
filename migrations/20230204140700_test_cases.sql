-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO outages (resource, city, district, start_date, end_date, source_url)
VALUES ('water','Limpopo','Ugadagada','2023-02-03 18:37:56', '2123-02-03 18:37:56', 'test entry'),
('water','Limpopo','Butumaputu','2023-02-03 18:37:56', '2123-02-03 18:37:56', 'test entry'),
('water','Limpopo','Ugadagada','2023-02-03 18:37:56', '2023-02-03 20:37:56', 'test entry')
;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

