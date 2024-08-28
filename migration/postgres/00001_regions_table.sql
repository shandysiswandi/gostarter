-- +goose Up
CREATE TABLE IF NOT EXISTS provinces (
  id CHAR(2) PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS cities (
  id CHAR(4) PRIMARY KEY,
  province_id CHAR(2) NOT NULL,
  type VARCHAR(10) NOT NULL DEFAULT '',
  name VARCHAR(255) NOT NULL,
  CONSTRAINT fk_province_id FOREIGN KEY (province_id) REFERENCES provinces(id)
);

CREATE TABLE IF NOT EXISTS districts (
  id CHAR(6) PRIMARY KEY,
  city_id CHAR(4) NOT NULL,
  name VARCHAR(255) NOT NULL,
  CONSTRAINT fk_city_id FOREIGN KEY (city_id) REFERENCES cities(id)
);

CREATE TABLE IF NOT EXISTS villages (
  id CHAR(10) PRIMARY KEY,
  district_id CHAR(6) NOT NULL,
  name VARCHAR(255) NOT NULL,
  postalcode CHAR(5) NOT NULL DEFAULT '',
  CONSTRAINT fk_district_id FOREIGN KEY (district_id) REFERENCES districts(id)
);

-- +goose Down
DROP TABLE IF EXISTS villages;
DROP TABLE IF EXISTS districts;
DROP TABLE IF EXISTS cities;
DROP TABLE IF EXISTS provinces;
