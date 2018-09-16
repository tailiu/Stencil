-- phpMyAdmin SQL Dump
-- version 4.8.0.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Sep 16, 2018 at 12:27 PM
-- Server version: 10.1.32-MariaDB
-- PHP Version: 5.6.36

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `stencil_storage`
--

-- --------------------------------------------------------

--
-- Table structure for table `apps`
--

CREATE TABLE `apps` (
  `PK` int(11) NOT NULL,
  `app_name` varchar(256) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `apps`
--

INSERT INTO `apps` (`PK`, `app_name`, `timestamp`) VALUES
(1, 'reddit', '2018-09-09 11:18:19'),
(2, 'twitter', '2018-09-09 11:18:19'),
(3, 'hacker news', '2018-09-09 11:18:19');

-- --------------------------------------------------------

--
-- Table structure for table `app_schemas`
--

CREATE TABLE `app_schemas` (
  `PK` int(11) NOT NULL,
  `table_id` int(11) NOT NULL,
  `column_name` varchar(256) NOT NULL,
  `data_type` int(11) NOT NULL,
  `constraints` varchar(512) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `app_schemas`
--

INSERT INTO `app_schemas` (`PK`, `table_id`, `column_name`, `data_type`, `constraints`, `timestamp`) VALUES
(1, 1, 'By', 1, NULL, '2018-09-10 07:25:20'),
(4, 1, 'Descendents', 1, NULL, '2018-09-12 10:26:07'),
(5, 1, 'Id', 1, NULL, '2018-09-12 10:27:17'),
(6, 1, 'Retrieved_on', 1, NULL, '2018-09-12 10:27:17'),
(7, 1, 'Score', 1, NULL, '2018-09-12 10:27:17'),
(8, 1, 'Time', 1, NULL, '2018-09-12 10:27:17'),
(9, 1, 'Title', 1, NULL, '2018-09-12 10:27:17'),
(10, 1, 'Type', 1, NULL, '2018-09-12 10:27:17'),
(11, 1, 'Url', 1, NULL, '2018-09-12 10:27:17'),
(12, 2, 'By', 1, NULL, '2018-09-12 10:30:16'),
(13, 2, 'Id', 1, NULL, '2018-09-12 10:30:16'),
(14, 2, 'Retrieved_on', 1, NULL, '2018-09-12 10:30:16'),
(15, 2, 'Time', 1, NULL, '2018-09-12 10:30:16'),
(16, 2, 'Kids', 1, NULL, '2018-09-12 10:30:16'),
(17, 2, 'Parent', 1, NULL, '2018-09-12 10:30:16'),
(18, 2, 'Text', 1, NULL, '2018-09-12 10:30:16'),
(19, 2, 'Type', 1, NULL, '2018-09-12 10:30:16'),
(20, 2, 'Author', 1, NULL, '2018-09-12 10:43:40');

-- --------------------------------------------------------

--
-- Table structure for table `app_tables`
--

CREATE TABLE `app_tables` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `table_name` varchar(256) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `app_tables`
--

INSERT INTO `app_tables` (`PK`, `app_id`, `table_name`, `timestamp`) VALUES
(1, 3, 'Story', '2018-09-16 10:09:44'),
(2, 3, 'Comment', '2018-09-16 10:09:44');

-- --------------------------------------------------------

--
-- Table structure for table `base_1`
--

CREATE TABLE `base_1` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) DEFAULT NULL,
  `row_id` text,
  `User` text,
  `Time` text,
  `URL` text,
  `Parent` text,
  `Text` text,
  `Id` text,
  `Retrieved_On` text,
  `Score` text,
  `Timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `base_2`
--

CREATE TABLE `base_2` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) DEFAULT NULL,
  `row_id` text,
  `Ups` text,
  `Num_Comments` text,
  `Name` text,
  `Profile_IMG` text,
  `Profile_Color` text,
  `Timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `base_table_attributes`
--

CREATE TABLE `base_table_attributes` (
  `PK` int(11) NOT NULL,
  `table_name` varchar(256) NOT NULL,
  `column_name` varchar(256) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `base_table_attributes`
--

INSERT INTO `base_table_attributes` (`PK`, `table_name`, `column_name`, `timestamp`) VALUES
(13, 'base_1', 'PK', '2018-09-12 11:04:38'),
(14, 'base_1', 'app_id', '2018-09-12 11:06:41'),
(15, 'base_1', 'row_id', '2018-09-12 11:06:41'),
(16, 'base_1', 'User', '2018-09-12 11:06:41'),
(17, 'base_1', 'Time', '2018-09-12 11:06:41'),
(18, 'base_1', 'URL', '2018-09-12 11:06:41'),
(19, 'base_1', 'Parent', '2018-09-12 11:06:41'),
(20, 'base_1', 'Text', '2018-09-12 11:06:41'),
(21, 'base_1', 'Id', '2018-09-12 11:06:41'),
(22, 'base_1', 'Retrieved_On', '2018-09-12 11:06:41'),
(23, 'base_1', 'Score', '2018-09-12 11:06:41'),
(24, 'base_1', 'Timestamp', '2018-09-12 11:06:41'),
(25, 'base_2', 'PK', '2018-09-12 11:08:07'),
(26, 'base_2', 'app_id', '2018-09-12 11:08:07'),
(27, 'base_2', 'row_id', '2018-09-12 11:08:07'),
(28, 'base_2', 'Ups', '2018-09-12 11:08:07'),
(29, 'base_2', 'Num_Comments', '2018-09-12 11:08:07'),
(30, 'base_2', 'Name', '2018-09-12 11:08:07'),
(31, 'base_2', 'Profile_IMG', '2018-09-12 11:08:07'),
(32, 'base_2', 'Profile_Color', '2018-09-12 11:08:07'),
(33, 'base_2', 'Timestamp', '2018-09-12 11:08:07');

-- --------------------------------------------------------

--
-- Table structure for table `data_types`
--

CREATE TABLE `data_types` (
  `PK` int(11) NOT NULL,
  `data_type` varchar(128) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `data_types`
--

INSERT INTO `data_types` (`PK`, `data_type`, `timestamp`) VALUES
(1, 'text', '2018-09-12 07:07:27'),
(2, 'int', '2018-09-12 07:07:27');

-- --------------------------------------------------------

--
-- Table structure for table `physical_mappings`
--

CREATE TABLE `physical_mappings` (
  `PK` int(11) NOT NULL,
  `logical_attribute` int(11) DEFAULT NULL,
  `physical_attribute` int(11) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `physical_mappings`
--

INSERT INTO `physical_mappings` (`PK`, `logical_attribute`, `physical_attribute`, `timestamp`) VALUES
(2, 1, 16, '2018-09-12 12:01:51'),
(4, 5, 21, '2018-09-12 12:01:51'),
(5, 6, 22, '2018-09-12 12:01:51'),
(6, 7, 23, '2018-09-12 12:01:51'),
(7, 8, 17, '2018-09-12 12:01:51'),
(8, 9, 20, '2018-09-12 12:01:51'),
(10, 11, 18, '2018-09-12 12:01:51'),
(11, 12, 16, '2018-09-12 12:01:51'),
(12, 13, 21, '2018-09-12 12:01:51'),
(13, 14, 22, '2018-09-12 12:01:51'),
(14, 15, 17, '2018-09-12 12:01:51'),
(16, 17, 19, '2018-09-12 12:01:51'),
(17, 18, 20, '2018-09-12 12:01:51');

-- --------------------------------------------------------

--
-- Table structure for table `schema_mappings`
--

CREATE TABLE `schema_mappings` (
  `PK` int(11) NOT NULL,
  `app1_attribute` int(11) NOT NULL,
  `app2_attribute` int(11) NOT NULL,
  `rules` varchar(512) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `schema_mappings`
--

INSERT INTO `schema_mappings` (`PK`, `app1_attribute`, `app2_attribute`, `rules`, `timestamp`) VALUES
(1, 1, 20, NULL, '2018-09-12 10:44:01');

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_1`
--

CREATE TABLE `supplementary_1` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `Descendents` text,
  `Type` text,
  `Timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_2`
--

CREATE TABLE `supplementary_2` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `Kids` text,
  `Retrieved_On` text,
  `Type` text,
  `Timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_3`
--

CREATE TABLE `supplementary_3` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `id_str` text,
  `truncated` text,
  `In_reply_to_status_id_str` text,
  `In_reply_to_user_id` text,
  `In_reply_to_user_id_str` text,
  `In_reply_to_screen_name` text,
  `geo` text,
  `coordinates` text,
  `lang` text,
  `user_mentions` text,
  `place` text,
  `contributors` text,
  `is_quote_status` text,
  `quote_count` text,
  `favorited` text,
  `retweet_count` text,
  `retweeted` text,
  `filter_level` text,
  `timestamp_ms` text,
  `hashtags` text,
  `symbols` text,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_4`
--

CREATE TABLE `supplementary_4` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `id_str` text,
  `name` text,
  `location` text,
  `url` text,
  `description` text,
  `translatior_type` text,
  `protected` text,
  `verified` text,
  `followers_count` text,
  `friends_count` text,
  `listed_count` text,
  `favourites_count` text,
  `statuses_count` text,
  `utc_offset` text,
  `time_zone` text,
  `geo_enabled` text,
  `lang` text,
  `contributors_enabled` text,
  `is_translator` text,
  `profile_background_image_url_https` text,
  ` profile_background_title` text,
  `profile_link_color` text,
  `profile_sidebar_border_color` text,
  `profile_sidebar_fill_color` text,
  `profile_text_color` text,
  `profile_use_background_image` text,
  `profile_image_url` text,
  `profile_image_url_https` text,
  `default_profile` text,
  `default_profile_image` text,
  `following` text,
  `following_request_sent` text,
  `notification` text,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_5`
--

CREATE TABLE `supplementary_5` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `link_karma` text,
  `Retrieved_On` text,
  `comment_karma` text,
  `profile_over_18` text,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_6`
--

CREATE TABLE `supplementary_6` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `Type` text,
  `author_url` text,
  `thumbernail_height` text,
  `thumbernail_url` text,
  `html` text,
  `Author_name` text,
  `provider_name` text,
  `Title` text,
  `provider_url` text,
  `version` text,
  `thumbnail_width` text,
  `width` text,
  `height` text,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_7`
--

CREATE TABLE `supplementary_7` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `height` text,
  `scrolling` text,
  `width` text,
  `content` text,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_8`
--

CREATE TABLE `supplementary_8` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `link_karma` text,
  `Image_id` text,
  `width` text,
  `height` text,
  `url` text,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_9`
--

CREATE TABLE `supplementary_9` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `Image_id` text,
  `Source_width` text,
  `Source_height` text,
  `Source_url` text,
  `Variants` text,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_10`
--

CREATE TABLE `supplementary_10` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `row_id` text,
  `archived` text,
  `link_flair_text` text,
  `Saved` text,
  `Thumb_nail` text,
  `link_flair_css_class` text,
  `spoiler` text,
  `edited` text,
  `domain` text,
  `hide_score` text,
  `contest_mode` text,
  `permalink` text,
  `distinguished` text,
  `subreddit_id` text,
  `name` text,
  `locked` text,
  `gilded` text,
  `subreddit` text,
  `over_18` text,
  `media_embed` text,
  `is_itself` text,
  ` author_flair_text` text,
  `stickied` text,
  `num_comments` text,
  `secure_media_embed` text,
  `quarantine` text,
  `preview` text,
  `post_hint` text,
  `Downs` text,
  `author_flair_css_class` text,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `supplementary_tables`
--

CREATE TABLE `supplementary_tables` (
  `PK` int(11) NOT NULL,
  `table_id` int(11) NOT NULL,
  `supplementary_table` varchar(256) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `supplementary_tables`
--

INSERT INTO `supplementary_tables` (`PK`, `table_id`, `supplementary_table`, `timestamp`) VALUES
(1, 1, 'supplementary_1', '2018-09-16 10:23:19'),
(2, 2, 'supplementary_2', '2018-09-16 10:23:19');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `apps`
--
ALTER TABLE `apps`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `app_schemas`
--
ALTER TABLE `app_schemas`
  ADD PRIMARY KEY (`PK`),
  ADD KEY `FK_data_type` (`data_type`),
  ADD KEY `FK_table_id` (`table_id`);

--
-- Indexes for table `app_tables`
--
ALTER TABLE `app_tables`
  ADD PRIMARY KEY (`PK`),
  ADD KEY `FK_app_id` (`app_id`);

--
-- Indexes for table `base_1`
--
ALTER TABLE `base_1`
  ADD PRIMARY KEY (`PK`),
  ADD KEY `Foreign Key` (`app_id`);

--
-- Indexes for table `base_2`
--
ALTER TABLE `base_2`
  ADD PRIMARY KEY (`PK`),
  ADD KEY `Foreign Key` (`app_id`);

--
-- Indexes for table `base_table_attributes`
--
ALTER TABLE `base_table_attributes`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `data_types`
--
ALTER TABLE `data_types`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `physical_mappings`
--
ALTER TABLE `physical_mappings`
  ADD PRIMARY KEY (`PK`),
  ADD KEY `FK_physical_attr` (`physical_attribute`),
  ADD KEY `FK_logical_attr` (`logical_attribute`);

--
-- Indexes for table `schema_mappings`
--
ALTER TABLE `schema_mappings`
  ADD PRIMARY KEY (`PK`),
  ADD KEY `FK_APP1_ATTR` (`app1_attribute`),
  ADD KEY `FK_APP2_ATTR` (`app2_attribute`);

--
-- Indexes for table `supplementary_1`
--
ALTER TABLE `supplementary_1`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_2`
--
ALTER TABLE `supplementary_2`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_3`
--
ALTER TABLE `supplementary_3`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_4`
--
ALTER TABLE `supplementary_4`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_5`
--
ALTER TABLE `supplementary_5`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_6`
--
ALTER TABLE `supplementary_6`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_7`
--
ALTER TABLE `supplementary_7`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_8`
--
ALTER TABLE `supplementary_8`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_9`
--
ALTER TABLE `supplementary_9`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_10`
--
ALTER TABLE `supplementary_10`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `supplementary_tables`
--
ALTER TABLE `supplementary_tables`
  ADD PRIMARY KEY (`PK`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `apps`
--
ALTER TABLE `apps`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `app_schemas`
--
ALTER TABLE `app_schemas`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

--
-- AUTO_INCREMENT for table `app_tables`
--
ALTER TABLE `app_tables`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `base_1`
--
ALTER TABLE `base_1`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `base_2`
--
ALTER TABLE `base_2`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `base_table_attributes`
--
ALTER TABLE `base_table_attributes`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=193;

--
-- AUTO_INCREMENT for table `data_types`
--
ALTER TABLE `data_types`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `physical_mappings`
--
ALTER TABLE `physical_mappings`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=19;

--
-- AUTO_INCREMENT for table `schema_mappings`
--
ALTER TABLE `schema_mappings`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `supplementary_1`
--
ALTER TABLE `supplementary_1`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_2`
--
ALTER TABLE `supplementary_2`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_3`
--
ALTER TABLE `supplementary_3`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_4`
--
ALTER TABLE `supplementary_4`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_5`
--
ALTER TABLE `supplementary_5`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_6`
--
ALTER TABLE `supplementary_6`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_7`
--
ALTER TABLE `supplementary_7`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_8`
--
ALTER TABLE `supplementary_8`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_9`
--
ALTER TABLE `supplementary_9`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_10`
--
ALTER TABLE `supplementary_10`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_tables`
--
ALTER TABLE `supplementary_tables`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `app_schemas`
--
ALTER TABLE `app_schemas`
  ADD CONSTRAINT `FK_data_type` FOREIGN KEY (`data_type`) REFERENCES `data_types` (`PK`),
  ADD CONSTRAINT `FK_table_id` FOREIGN KEY (`table_id`) REFERENCES `app_tables` (`PK`);

--
-- Constraints for table `app_tables`
--
ALTER TABLE `app_tables`
  ADD CONSTRAINT `FK_app_id` FOREIGN KEY (`app_id`) REFERENCES `apps` (`PK`);

--
-- Constraints for table `base_1`
--
ALTER TABLE `base_1`
  ADD CONSTRAINT `Foreign Key` FOREIGN KEY (`app_id`) REFERENCES `apps` (`PK`);

--
-- Constraints for table `physical_mappings`
--
ALTER TABLE `physical_mappings`
  ADD CONSTRAINT `FK_logical_attr` FOREIGN KEY (`logical_attribute`) REFERENCES `app_schemas` (`PK`),
  ADD CONSTRAINT `FK_physical_attr` FOREIGN KEY (`physical_attribute`) REFERENCES `base_table_attributes` (`PK`);

--
-- Constraints for table `schema_mappings`
--
ALTER TABLE `schema_mappings`
  ADD CONSTRAINT `FK_APP1_ATTR` FOREIGN KEY (`app1_attribute`) REFERENCES `app_schemas` (`PK`),
  ADD CONSTRAINT `FK_APP2_ATTR` FOREIGN KEY (`app2_attribute`) REFERENCES `app_schemas` (`PK`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
