SET @OLD_TIME_ZONE=@@TIME_ZONE;
SET TIME_ZONE='+00:00';

CREATE DATABASE `db_apiserver`;
USE `db_apiserver`;

DROP TABLE IF EXISTS `tb_users`;
CREATE TABLE `tb_users`(
  id bigint(20) unsigned not null auto_increment,
  username varchar(255) not null,
  password varchar(255) not null,
  createdAt timestamp null DEFAULT null,
  updatedAt timestamp null default null,
  deletedAt timestamp null default null,
  primary key (id),
  UNIQUE key username (username),
  KEY idx_tb_users_deletedAt (deleteAt)
) AUTO_INCREMENT=2;

LOCK TABLES `tb_users` WRITE;
INSERT INTO `tb_users` VALUES (0,'admin','$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG','2018-08-13 10:00:00','2018-08-13 10:00:00',null);
UNLOCK tables;