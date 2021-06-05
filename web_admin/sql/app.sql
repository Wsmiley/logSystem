USE logadmin;

DROP TABLE IF EXISTS `tbl_log_info`;
CREATE TABLE `tbl_log_info` (
  `log_id` int(32) NOT NULL AUTO_INCREMENT,
  `app_id` int(32) NOT NULL ,
  `app_name` varchar(255) NOT NULL,
  `topic` varchar(255) NOT NULL,
  `log_path` varchar(255) NOT NULL,
  `create_time`datetime   default NULL,
  `ip` varchar(255) NOT NULL,
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `tbl_app_info`;
CREATE TABLE `tbl_app_info` (
  `app_id` int(32) NOT NULL AUTO_INCREMENT,
  `app_name` varchar(255) NOT NULL,
  `app_type` varchar(255) NOT NULL,
  `develop_path` varchar(255) NOT NULL,
  `create_time`datetime    default NULL,
  PRIMARY KEY (`app_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `tbl_app_ip`;
CREATE TABLE `tbl_app_ip` (
   `ip_id` int(32) NOT NULL AUTO_INCREMENT,
  `app_id` int(32) NOT NULL,
  
  `ip` varchar(255) NOT NULL,
  `create_time` datetime  default NULL,
  PRIMARY KEY (`ip_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `tbl_admin`;
CREATE TABLE `tbl_admin` (
  `admin_id` int(32) UNSIGNED AUTO_INCREMENT,
  `username` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `status` INT(4),
  `create_time` datetime  default NULL ,
  PRIMARY KEY (`admin_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
