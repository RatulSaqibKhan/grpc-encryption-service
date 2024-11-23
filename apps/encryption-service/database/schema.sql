CREATE DATABASE `demodb`;
USE `demodb`;

DROP TABLE IF EXISTS `encryption`;
CREATE TABLE `encryption` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `plaintext` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `encrypted` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;