CREATE TABLE `user_feeds`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`  timestamp NULL DEFAULT NULL,
    `telegram_id` bigint    NOT NULL,
    `feed_id`     bigint unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `feed_telegram` (`feed_id`,`telegram_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
