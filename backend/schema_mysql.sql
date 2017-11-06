CREATE TABLE `jobs` (
  `id` char(36) NOT NULL,
  `command` varchar(255) NOT NULL,
  `args` text NOT NULL,
  `flags` text NOT NULL,
  `user_id` char(36) DEFAULT NULL,
  `state` varchar(20) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `started_at` datetime DEFAULT NULL,
  `finished_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `sessions` (
  `id` char(36) NOT NULL,
  `user_id` char(36) NOT NULL,
  `created_at` datetime NOT NULL,
  `expires_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `users` (
  `id` char(36) NOT NULL,
  `github_id` int(11) unsigned NOT NULL,
  `github_login` varchar(255) NOT NULL DEFAULT '',
  `github_name` varchar(255) NOT NULL DEFAULT '',
  `github_picture` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
