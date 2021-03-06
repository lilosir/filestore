CREATE TABLE `tbl_file` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT 'file hash',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'file name',
    `file_size` bigint(20) DEFAULT '0' COMMENT 'file size',
    `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT 'file address',
    `create_at` datetime default NOW() COMMENT 'create time',
    `update_at` datetime default NOW() on update current_timestamp() COMMENT 'update time',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT 'status(available/disabled/deleted, etc)',
    `ext1` int(11) DEFAULT '0' COMMENT 'extension 1',
    `ext2` text COMMENT 'extension 2',
    PRIMARY KEY (`id`),
    UNIQUE KEY `id_file_hash` (`file_sha1`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tbl_user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'user name',
    `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT 'user password',
    `email` varchar(64) DEFAULT '' COMMENT 'user email',
    `phone` varchar(64) DEFAULT '' COMMENT 'user phone number',
    `email_validated` tinyint(1) DEFAULT 0 COMMENT 'if the email is valided',
    `phone_validated` tinyint(1) DEFAULT 0 COMMENT 'if the phone number is valided',
    `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'register time',
    `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last active time',
    `profile` text COMMENT 'user profile',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT 'status(available/disabled/deleted, etc)',
    PRIMARY KEY (`id`),
    KEY `idx_status`(`status`)
)ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `tbl_user_token` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'user name',
    `user_token` char(40) NOT NULL DEFAULT '' COMMENT 'authentication token',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username`(`user_name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `tbl_user_file` (
    `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'user name',
    `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT 'file shah',
    `file_size` bigint(20) DEFAULT '0' COMMENT 'file size',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'file name',
    `upload_at` datetime DEFAULT CURRENT_TIMESTAMP comment 'uploadt time',
    `last_update` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last updated time',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT 'status(available/disabled/deleted, etc)',
    UNIQUE KEY `idx_user_file` (`user_name`, `file_sha1`),
    KEY `idx_status`(`status`),
    KEY `idx_user_id`(`user_name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;