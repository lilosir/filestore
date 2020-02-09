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