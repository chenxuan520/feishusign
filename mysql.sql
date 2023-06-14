CREATE TABLE `sign` (
  `user_id` varchar(255) NOT NULL,
  `user_name` varchar(255) DEFAULT NULL,
  `status` tinyint(4) NOT NULL,
  `meeting_id` bigint(20) NOT NULL,
  `create_time` bigint(20) NOT NULL,
  PRIMARY KEY (`user_id`,`meeting_id`),
  KEY `idx_meeting_id` (`meeting_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
