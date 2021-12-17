BEGIN;

CREATE DATABASE IF NOT EXISTS lake_test;

CREATE USER IF NOT EXISTS 'merico'@'localhost' IDENTIFIED BY 'merico';
GRANT ALL PRIVILEGES ON *.* TO 'merico'@'%';

USE lake_test;

CREATE TABLE `schema_migrations` (
  `version` bigint NOT NULL DEFAULT 1,
  `dirty` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SELECT now();

COMMIT;
