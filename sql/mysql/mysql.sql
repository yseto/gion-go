SET foreign_key_checks=0;

CREATE TABLE `users` (
  `id` integer(11) NOT NULL auto_increment,
  `name` varchar(32) NULL DEFAULT NULL,
  `digest` varchar(200) NULL,
  `numentry` integer(11) NOT NULL DEFAULT 0,
  `nopinlist` integer(11) NOT NULL DEFAULT 0,
  `numsubstr` integer(11) NOT NULL DEFAULT 0,
  `autoseen` boolean NOT NULL DEFAULT 0,
  `last_login` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `category` (
  `id` integer(11) NOT NULL auto_increment,
  `user_id` integer(11) NOT NULL,
  `name` varchar(60) NOT NULL,
  INDEX `user_id` (`user_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `category_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `feed` (
  `id` integer(11) NOT NULL auto_increment,
  `url` text NOT NULL,
  `siteurl` text NOT NULL,
  `title` varchar(200) NOT NULL,
  `http_status` varchar(3) NOT NULL,
  `pubdate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `term` varchar(1) NOT NULL DEFAULT '1',
  `cache` text NOT NULL,
  `next_serial` integer(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `subscription` (
  `id` integer(11) NOT NULL auto_increment,
  `category_id` integer(11) NOT NULL,
  `feed_id` integer(11) NOT NULL,
  `user_id` integer(11) NOT NULL,
  INDEX `category_id` (`category_id`),
  INDEX `feed_id` (`feed_id`),
  INDEX (`user_id`),
  PRIMARY KEY (`id`),
  CONSTRAINT `subscription_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `subscription_ibfk_2` FOREIGN KEY (`feed_id`) REFERENCES `feed` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `subscription_ibfk_3` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `entry` (
  `serial` integer(11) NOT NULL,
  `pubdate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `readflag` tinyint(4) NOT NULL,
  `subscription_id` integer(11) NOT NULL,
  `feed_id` integer(11) NOT NULL,
  `user_id` integer(11) NOT NULL,
  INDEX `subscription_id` (`subscription_id`),
  INDEX (`user_id`),
  UNIQUE `serial_2` (`serial`, `feed_id`, `user_id`),
  CONSTRAINT `entry_ibfk_1` FOREIGN KEY (`subscription_id`) REFERENCES `subscription` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `entry_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE `story` (
  `feed_id` integer(11) NOT NULL,
  `serial` integer(11) NOT NULL,
  `title` varchar(80) NOT NULL,
  `description` tinytext NOT NULL,
  `url` text NOT NULL,
  PRIMARY KEY (`feed_id`, `serial`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

SET foreign_key_checks=1;

