-- +goose Up
-- SQL in this section is executed when the migration is applied.

delete from districts where city='Datça';

INSERT INTO districts  (city, district)
  	VALUES ('Fethiye','Çamköy'),
  	('Fethiye','Karacaören'),
	('Fethiye','Sögütlü'),
	('Fethiye','Yanıklar'),
	('Bodrum','Kumbahçe'),
	('Bodrum','Çiftlik'),
	('Bodrum','Çömlekçi'),
	('Bodrum','Gölköy'),
	('Bodrum','Kızıkağaç'),
	('Bodrum','Müskebi'),
	('Bodrum','Tepecik Karaova'),
	('Bodrum','Tepecik'),
	('Bodrum','Türkbükü'),
	('Bodrum','Yahşi'),
	('Bodrum','Yalıkavak'),
	('Bodrum','Yeniköy Karaova'),
	('Bodrum','Yeniköy'),
	('Dalaman','Gürleyk'),
	('Dalaman','Hürriyet'),
	('Dalaman','Merkez'),
	('Dalaman','Narlı'),
	('Datça','Cumalı'),
	('Datça','Datça'),
	('Datça','Emecik'),
	('Datça','Hızırşah'),
	('Datça','İskele'),
	('Datça','Karaköy'),
	('Datça','Kızlan'),
	('Datça','Mesudiye'),
	('Datça','Reşadiye'),
	('Datça','Sındı'),
	('Datça','Yaka'),
	('Datça','Yazı'),
	('Kavaklıdere','Cumhuriyet'),
	('Köyceğiz','Kavakarasi'),
	('Menteşe','Çatakbağyaka'),
	('Menteşe','Dokuzçam'),
	('Menteşe','İkizce'),
	('Menteşe','Kuzluk'),
	('Menteşe','Şenyayla'),
	('Menteşe','Sungur'),
	('Menteşe','Yeniköy Yerkesik'),
	('Menteşe','Yerkesik'),
	('Milas','Beyciler'),
	('Milas','Çiftlik'),
	('Milas','Dibekdere'),
	('Milas','Karapınar'),
	('Milas','Kılavuz'),
	('Milas','Kısırlar'),
	('Milas','Küçükdibekdere'),
	('Milas','Sek'),
	('Milas','Yusufca'),
	('Ortaca','Fevziye'),
	('Seydikemer','Belen'),
	('Seydikemer','Boğaziçi'),
	('Seydikemer','Cumhuriyet'),
	('Seydikemer','Dodurga'),
	('Seydikemer','Gerişburnu'),
	('Seydikemer','İzzettinköy'),
	('Seydikemer','Kabaağaç'),
	('Seydikemer','Kıncılar'),
	('Seydikemer','Menekşe'),
	('Seydikemer','Minare'),
	('Seydikemer','Ören'),
	('Seydikemer','Sarıyer'),
	('Seydikemer','Söğütlüdere'),
	('Seydikemer','Yayla Eldirek'),
	('Seydikemer','Yayla Gökben'),
	('Seydikemer','Yayla Karaçulha'),
	('Seydikemer','Yaylapatlangıç'),
	('Ula','Ataköy'),
	('Ula','Çörüş'),
	('Ula','Gökçe'),
	('Ula','Portakallık'),
	('Ula','Yeşilçam'),
	('Yatağan','Bencik'),
	('Yatağan','Çamlıca'),
	('Yatağan','Kozağaç'),
	('Yatağan','Mesken');

ALTER TABLE districts 
ADD column IF NOT EXISTS city_normalized TEXT,
ADD COLUMN IF NOT exists district_normalized TEXT;

CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

UPDATE districts 
SET city_normalized = lower (normalize(regexp_replace(unaccent(city), '[^a-zA-Z0-9\s]+', ' ', 'g'), NFD)),
district_normalized = lower (normalize(regexp_replace(unaccent(district), '[^a-zA-Z0-9\s]+', ' ', 'g'),NFD));

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP EXTENSION IF EXISTS "pg_trgm",
DROP EXTENSION IF EXISTS "unaccent";
DROP column IF EXISTS city_normalized,
DROP column IF EXISTS district_normalized,