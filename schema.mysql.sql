DROP DATABASE `gocontacts`;

CREATE DATABASE `gocontacts`;
USE `gocontacts`;

CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `logged_out_at` timestamp NULL DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `role` varchar(255) DEFAULT NULL,
  `state` varchar(16) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=latin1;

CREATE TABLE `sign_up_confirmations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `nonce` varchar(32) NOT NULL,
  `user_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `nonce` (`nonce`),
  FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
  KEY `idx_contacts_user_id` (`user_id`),
  KEY `idx_contacts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

CREATE TABLE `contacts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `user_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
  KEY `idx_contacts_user_id` (`user_id`),
  KEY `idx_contacts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

CREATE TABLE `card_transactions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `datetime` timestamp DEFAULT CURRENT_TIMESTAMP,
  `amount` BIGINT NOT NULL,
  `currency_scale` TINYINT NOT NULL,
  `currency_code` varchar(255) NOT NULL,
  `reference` varchar(255) NOT NULL,
  `merchant_name` varchar(255) NOT NULL,
  `merchant_city` varchar(255) NOT NULL,
  `merchant_country_code` varchar(255) NOT NULL,
  `merchant_country_name` varchar(255) NOT NULL,
  `merchant_category_code` varchar(255) NOT NULL,
  `merchant_category_name` varchar(255) NOT NULL,
  `user_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
  KEY `idx_contacts_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;