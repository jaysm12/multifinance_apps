CREATE DATABASE `multifinance-apps`;

USE `multifinance-apps`;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(50) NOT NULL,
  `password` longtext NOT NULL,
  `fullname` longtext,
  `email` longtext,
  `is_verified` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `user_kycs`;
CREATE TABLE `user_kycs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `nik` varchar(16) NOT NULL,
  `legal_name` longtext,
  `birth_date` longtext,
  `birth_address` longtext,
  `salary_amount` longtext,
  `photo_id_url` longtext,
  `photo_selfie_url` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_user_kycs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `credit_options`;
CREATE TABLE `credit_options` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `default_amount` double DEFAULT NULL,
  `current_amount` double DEFAULT NULL,
  `tenor` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_credit_options_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


LOCK TABLES `credit_options` WRITE;
UNLOCK TABLES;


DROP TABLE IF EXISTS `installments`;
CREATE TABLE `installments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `credit_option_id` bigint unsigned DEFAULT NULL,
  `contract_id` longtext,
  `asset_name` longtext,
  `otr_amount` double DEFAULT NULL,
  `admin_fee` bigint DEFAULT NULL,
  `total_interest` double DEFAULT NULL,
  `monthly_amount` double DEFAULT NULL,
  `total_installment_amount` double DEFAULT NULL,
  `interest_rate` double DEFAULT NULL,
  `interest_per_month` double DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `overdue_amount` bigint DEFAULT NULL,
  `overdue_days` bigint DEFAULT NULL,
  `tenor` bigint DEFAULT NULL,
  `remaining_amount` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_installments_deleted_at` (`deleted_at`),
  KEY `fk_credit_options_installments` (`credit_option_id`),
  CONSTRAINT `fk_credit_options_installments` FOREIGN KEY (`credit_option_id`) REFERENCES `credit_options` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `installment_payment_histories`;
CREATE TABLE `installment_payment_histories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `installment_id` bigint unsigned DEFAULT NULL,
  `contract_id` longtext,
  `installment_number` bigint DEFAULT NULL,
  `payment_date` datetime(3) DEFAULT NULL,
  `paid_amount` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_installment_payment_histories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;