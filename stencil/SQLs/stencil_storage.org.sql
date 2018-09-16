-- phpMyAdmin SQL Dump
-- version 4.8.0.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Sep 11, 2018 at 03:03 PM
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
-- Table structure for table `app_mappings`
--

CREATE TABLE `app_mappings` (
  `PK` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  `table_id` int(11) DEFAULT NULL,
  `column_name` varchar(256) NOT NULL,
  `data_type` varchar(512) NOT NULL,
  `mapping` varchar(256) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `app_mappings`
--

INSERT INTO `app_mappings` (`PK`, `app_id`, `table_id`, `column_name`, `data_type`, `mapping`, `timestamp`) VALUES
(1, 3, 1, 'By', 'Text', 'base_1.User', '2018-09-09 11:21:58'),
(2, 3, 1, 'Descendents', 'Text', 'supplementary_1.Descendents', '2018-09-09 11:24:22'),
(3, 3, 1, 'Id', 'Text', 'base_1.Id', '2018-09-09 11:24:22'),
(4, 3, 1, 'Retrieved_on', 'Text', 'base_1.Retrieved_on', '2018-09-09 11:24:22'),
(5, 3, 1, 'Score', 'Text', 'base_1.Score', '2018-09-09 11:24:22'),
(6, 3, 1, 'Time', 'Text', 'base_1.Time', '2018-09-09 11:24:22'),
(7, 3, 1, 'Title', 'Text', 'base_1.Text', '2018-09-09 11:24:22'),
(8, 3, 1, 'Type', 'Text', 'supplementary_1.Type', '2018-09-09 11:24:22'),
(9, 3, 1, 'Url', 'Text', 'base_1.Url', '2018-09-09 11:24:22'),
(10, 3, 2, 'By', 'Text', 'base_1.User', '2018-09-09 11:24:22'),
(11, 3, 2, 'Id', 'Text', 'base_1.Id', '2018-09-09 11:24:22'),
(12, 3, 2, 'Retrieved_on', 'Text', 'base_1.Retrieved_on', '2018-09-09 11:24:22'),
(13, 3, 2, 'Time', 'Text', 'base_1.Time', '2018-09-09 11:24:22'),
(14, 3, 2, 'Kids', 'Text', 'supplementary_2.Kids', '2018-09-09 11:24:22'),
(15, 3, 2, 'Parent', 'Text', 'base_1.Parent', '2018-09-09 11:24:22'),
(16, 3, 2, 'Text', 'Text', 'base_1.Text', '2018-09-09 11:24:22'),
(17, 3, 2, 'Type', 'Text', 'supplementary_2.Type', '2018-09-09 11:24:22');

-- --------------------------------------------------------

--
-- Table structure for table `app_tables`
--

CREATE TABLE `app_tables` (
  `PK` int(11) NOT NULL,
  `table_name` varchar(256) NOT NULL,
  `app_id` int(11) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `app_tables`
--

INSERT INTO `app_tables` (`PK`, `table_name`, `app_id`, `timestamp`) VALUES
(1, 'Story', 3, '2018-09-10 07:25:20'),
(2, 'Comment', 3, '2018-09-10 07:25:20');

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

--
-- Indexes for dumped tables
--

--
-- Indexes for table `apps`
--
ALTER TABLE `apps`
  ADD PRIMARY KEY (`PK`);

--
-- Indexes for table `app_mappings`
--
ALTER TABLE `app_mappings`
  ADD PRIMARY KEY (`PK`),
  ADD KEY `foreign_key_table` (`table_id`),
  ADD KEY `foreign_key_app` (`app_id`) USING BTREE;

--
-- Indexes for table `app_tables`
--
ALTER TABLE `app_tables`
  ADD PRIMARY KEY (`PK`),
  ADD KEY `foreign_key_apps` (`app_id`);

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
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `apps`
--
ALTER TABLE `apps`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `app_mappings`
--
ALTER TABLE `app_mappings`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=18;

--
-- AUTO_INCREMENT for table `app_tables`
--
ALTER TABLE `app_tables`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `base_1`
--
ALTER TABLE `base_1`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=204;

--
-- AUTO_INCREMENT for table `base_2`
--
ALTER TABLE `base_2`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `supplementary_1`
--
ALTER TABLE `supplementary_1`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=26;

--
-- AUTO_INCREMENT for table `supplementary_2`
--
ALTER TABLE `supplementary_2`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=179;

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
-- Constraints for dumped tables
--

--
-- Constraints for table `app_mappings`
--
ALTER TABLE `app_mappings`
  ADD CONSTRAINT `foreign_key_app` FOREIGN KEY (`app_id`) REFERENCES `apps` (`PK`),
  ADD CONSTRAINT `foreign_key_table` FOREIGN KEY (`table_id`) REFERENCES `app_tables` (`PK`);

--
-- Constraints for table `app_tables`
--
ALTER TABLE `app_tables`
  ADD CONSTRAINT `foreign_key_apps` FOREIGN KEY (`app_id`) REFERENCES `apps` (`PK`);

--
-- Constraints for table `base_1`
--
ALTER TABLE `base_1`
  ADD CONSTRAINT `Foreign Key` FOREIGN KEY (`app_id`) REFERENCES `apps` (`PK`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
