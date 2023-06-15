CREATE TABLE `sign` (
  `user_id` varchar(255) NOT NULL,
  `meeting_id` varchar(255) NOT NULL,
  `user_name` varchar(255) DEFAULT NULL,
  `status` tinyint(4) NOT NULL,
  `create_time` bigint(20) NOT NULL,
  PRIMARY KEY (`user_id`,`meeting_id`),
  KEY `idx_meeting_id` (`meeting_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `meeting` (
  `meeting_id` VARCHAR(255),
  `originator_id` VARCHAR(255),
  `year` INT,
  `month` INT,
  `day` INT,
  `create_time` BIGINT,
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

