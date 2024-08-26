-- +goose Up
CREATE TABLE IF NOT EXISTS `provinces` (
  `id` char(2) NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `cities` (
  `id` char(4) NOT NULL,
  `province_id` char(2) NOT NULL,
  `type` char(10) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `province_id` (`province_id`),
  CONSTRAINT `cities_ibfk_1` FOREIGN KEY (`province_id`) REFERENCES `provinces` (`id`)
);

CREATE TABLE IF NOT EXISTS `districts` (
  `id` char(6) NOT NULL,
  `city_id` char(4) NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `city_id` (`city_id`),
  CONSTRAINT `districts_ibfk_1` FOREIGN KEY (`city_id`) REFERENCES `cities` (`id`)
);

CREATE TABLE IF NOT EXISTS `villages` (
  `id` char(10) NOT NULL,
  `district_id` char(6) NOT NULL,
  `name` varchar(255) NOT NULL,
  `postalcode` char(5) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `district_id` (`district_id`),
  CONSTRAINT `villages_ibfk_1` FOREIGN KEY (`district_id`) REFERENCES `districts` (`id`)
);

-- +goose Down
DROP TABLE IF EXISTS villages;
DROP TABLE IF EXISTS districts;
DROP TABLE IF EXISTS cities;
DROP TABLE IF EXISTS provinces;
