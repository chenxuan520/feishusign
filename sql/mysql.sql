CREATE TABLE `sign` (
  `user_id` varchar(255) NOT NULL,
  `meeting_id` varchar(255) NOT NULL,
  `user_name` varchar(255) DEFAULT NULL,
  `status` tinyint(4) NOT NULL,
  `create_time` bigint(20) NOT NULL,
  PRIMARY KEY (`user_id`,`meeting_id`),
  KEY `idx_meeting_id` (`meeting_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE meeting (
    meeting_id VARCHAR(255) NOT NULL,
    originator_id VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    month INT NOT NULL,
    day INT NOT NULL,
    create_time BIGINT NOT NULL,
    PRIMARY KEY (meeting_id),
    UNIQUE INDEX meeting_id_UNIQUE (meeting_id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
