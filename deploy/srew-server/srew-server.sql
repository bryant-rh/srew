DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL DEFAULT '' unique,
  `password` varchar(100) NOT NULL DEFAULT '',
	`created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `plugin`;
CREATE TABLE `plugin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `plugin_id` varchar(50)  NOT NULL DEFAULT '0' unique,
  `plugin_name` varchar(50) NOT NULL DEFAULT '',
  `latest_version` varchar(50) NOT NULL DEFAULT '',
	`created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `detail`;
CREATE TABLE `detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `plugin_id` varchar(50) NOT NULL DEFAULT '0',
  `plugin_name` varchar(50) NOT NULL DEFAULT '',
  `version` varchar(50) NOT NULL DEFAULT '',
	`homepage` varchar(255) NOT NULL DEFAULT '',
	`shortDescription` varchar(255) NOT NULL DEFAULT '',
	`description` varchar(255) NOT NULL DEFAULT '',
	`caveats` varchar(255) NOT NULL DEFAULT '',
	`platforms` varchar(2048) NOT NULL DEFAULT '',
	`created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
