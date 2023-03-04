CREATE TABLE `feeds`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at`  timestamp                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  timestamp                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`  timestamp NULL DEFAULT NULL,
    `address`     varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `address_md5` char(32)                                                       NOT NULL,
    `latest_post` varchar(2048)                                                  NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `feed_address_md5_IDX` (`address_md5`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
