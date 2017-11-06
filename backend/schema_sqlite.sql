CREATE TABLE `jobs` (
  `id` CHAR(36) NOT NULL PRIMARY KEY,
  `command` VARCHAR(255) NOT NULL,
  `args` TEXT NOT NULL,
  `flags` TEXT NOT NULL,
  `user_id` CHAR(36) DEFAULT NULL,
  `state` VARCHAR(20) NOT NULL,
  `created_at` DATETIME DEFAULT NULL,
  `started_at` DATETIME DEFAULT NULL,
  `finished_at` DATETIME DEFAULT NULL
);
-- PRAGMA table_info(jobs);

CREATE TABLE `sessions` (
  `id` CHAR(36) NOT NULL PRIMARY KEY,
  `user_id` CHAR(36) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `expires_at` DATETIME NOT NULL
);
-- PRAGMA table_info(sessions);

CREATE TABLE `users` (
  `id` CHAR(36) NOT NULL PRIMARY KEY,
  `github_id` UNSIGNED BIGINT NOT NULL,
  `github_login` VARCHAR(255) NOT NULL DEFAULT '',
  `github_name` VARCHAR(255) NOT NULL DEFAULT '',
  `github_picture` VARCHAR(255) NOT NULL DEFAULT ''
);
-- PRAGMA table_info(users);
