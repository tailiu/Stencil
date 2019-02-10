-- phpMyAdmin SQL Dump
-- version 4.8.0.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Sep 12, 2018 at 02:51 PM
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
  `app_id` int(11) NOT NULL,
  `table_name` varchar(256) NOT NULL,
  `column_name` varchar(256) NOT NULL,
  `data_type` int(11) NOT NULL,
  `constraints` varchar(512) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `app_schemas`
--

INSERT INTO `app_schemas` (`PK`, `app_id`, `table_name`, `column_name`, `data_type`, `constraints`, `timestamp`) VALUES
(1, 3, 'Story', 'By', 1, NULL, '2018-09-10 07:25:20'),
(4, 3, 'Story', 'Descendents', 1, NULL, '2018-09-12 10:26:07'),
(5, 3, 'Story', 'Id', 1, NULL, '2018-09-12 10:27:17'),
(6, 3, 'Story', 'Retrieved_on', 1, NULL, '2018-09-12 10:27:17'),
(7, 3, 'Story', 'Score', 1, NULL, '2018-09-12 10:27:17'),
(8, 3, 'Story', 'Time', 1, NULL, '2018-09-12 10:27:17'),
(9, 3, 'Story', 'Title', 1, NULL, '2018-09-12 10:27:17'),
(10, 3, 'Story', 'Type', 1, NULL, '2018-09-12 10:27:17'),
(11, 3, 'Story', 'Url', 1, NULL, '2018-09-12 10:27:17'),
(12, 3, 'Comment', 'By', 1, NULL, '2018-09-12 10:30:16'),
(13, 3, 'Comment', 'Id', 1, NULL, '2018-09-12 10:30:16'),
(14, 3, 'Comment', 'Retrieved_on', 1, NULL, '2018-09-12 10:30:16'),
(15, 3, 'Comment', 'Time', 1, NULL, '2018-09-12 10:30:16'),
(16, 3, 'Comment', 'Kids', 1, NULL, '2018-09-12 10:30:16'),
(17, 3, 'Comment', 'Parent', 1, NULL, '2018-09-12 10:30:16'),
(18, 3, 'Comment', 'Text', 1, NULL, '2018-09-12 10:30:16'),
(19, 3, 'Comment', 'Type', 1, NULL, '2018-09-12 10:30:16'),
(20, 1, 'Submission', 'Author', 1, NULL, '2018-09-12 10:43:40');

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
(3, 4, 37, '2018-09-12 12:01:51'),
(4, 5, 21, '2018-09-12 12:01:51'),
(5, 6, 22, '2018-09-12 12:01:51'),
(6, 7, 23, '2018-09-12 12:01:51'),
(7, 8, 17, '2018-09-12 12:01:51'),
(8, 9, 20, '2018-09-12 12:01:51'),
(9, 10, 38, '2018-09-12 12:01:51'),
(10, 11, 18, '2018-09-12 12:01:51'),
(11, 12, 16, '2018-09-12 12:01:51'),
(12, 13, 21, '2018-09-12 12:01:51'),
(13, 14, 22, '2018-09-12 12:01:51'),
(14, 15, 17, '2018-09-12 12:01:51'),
(15, 16, 43, '2018-09-12 12:01:51'),
(16, 17, 19, '2018-09-12 12:01:51'),
(17, 18, 20, '2018-09-12 12:01:51'),
(18, 19, 45, '2018-09-12 12:01:51');

-- --------------------------------------------------------

--
-- Table structure for table `physical_schemas`
--

CREATE TABLE `physical_schemas` (
  `PK` int(11) NOT NULL,
  `table_name` varchar(256) NOT NULL,
  `column_name` varchar(256) NOT NULL,
  `type` varchar(64) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `physical_schemas`
--

INSERT INTO `physical_schemas` (`PK`, `table_name`, `column_name`, `type`, `timestamp`) VALUES
(13, 'base_1', 'PK', 'b', '2018-09-12 11:04:38'),
(14, 'base_1', 'app_id', 'b', '2018-09-12 11:06:41'),
(15, 'base_1', 'row_id', 'b', '2018-09-12 11:06:41'),
(16, 'base_1', 'User', 'b', '2018-09-12 11:06:41'),
(17, 'base_1', 'Time', 'b', '2018-09-12 11:06:41'),
(18, 'base_1', 'URL', 'b', '2018-09-12 11:06:41'),
(19, 'base_1', 'Parent', 'b', '2018-09-12 11:06:41'),
(20, 'base_1', 'Text', 'b', '2018-09-12 11:06:41'),
(21, 'base_1', 'Id', 'b', '2018-09-12 11:06:41'),
(22, 'base_1', 'Retrieved_On', 'b', '2018-09-12 11:06:41'),
(23, 'base_1', 'Score', 'b', '2018-09-12 11:06:41'),
(24, 'base_1', 'Timestamp', 'b', '2018-09-12 11:06:41'),
(25, 'base_2', 'PK', 'b', '2018-09-12 11:08:07'),
(26, 'base_2', 'app_id', 'b', '2018-09-12 11:08:07'),
(27, 'base_2', 'row_id', 'b', '2018-09-12 11:08:07'),
(28, 'base_2', 'Ups', 'b', '2018-09-12 11:08:07'),
(29, 'base_2', 'Num_Comments', 'b', '2018-09-12 11:08:07'),
(30, 'base_2', 'Name', 'b', '2018-09-12 11:08:07'),
(31, 'base_2', 'Profile_IMG', 'b', '2018-09-12 11:08:07'),
(32, 'base_2', 'Profile_Color', 'b', '2018-09-12 11:08:07'),
(33, 'base_2', 'Timestamp', 'b', '2018-09-12 11:08:07'),
(34, 'supplementary_1', 'PK', 's', '2018-09-12 11:10:10'),
(35, 'supplementary_1', 'app_id', 's', '2018-09-12 11:10:10'),
(36, 'supplementary_1', 'row_id', 's', '2018-09-12 11:10:10'),
(37, 'supplementary_1', 'Descendents', 's', '2018-09-12 11:10:10'),
(38, 'supplementary_1', 'Type', 's', '2018-09-12 11:10:10'),
(39, 'supplementary_1', 'Timestamp', 's', '2018-09-12 11:10:10'),
(40, 'supplementary_2', 'PK', 's', '2018-09-12 11:11:30'),
(41, 'supplementary_2', 'app_id', 's', '2018-09-12 11:11:30'),
(42, 'supplementary_2', 'row_id', 's', '2018-09-12 11:11:30'),
(43, 'supplementary_2', 'Kids', 's', '2018-09-12 11:11:30'),
(44, 'supplementary_2', 'Retrieved_On', 's', '2018-09-12 11:11:30'),
(45, 'supplementary_2', 'Type', 's', '2018-09-12 11:11:30'),
(46, 'supplementary_2', 'Timestamp', 's', '2018-09-12 11:11:30'),
(47, 'supplementary_3', 'PK', 's', '2018-09-12 11:20:20'),
(48, 'supplementary_3', 'app_id', 's', '2018-09-12 11:20:20'),
(49, 'supplementary_3', 'row_id', 's', '2018-09-12 11:20:20'),
(50, 'supplementary_3', 'id_str', 's', '2018-09-12 11:20:20'),
(51, 'supplementary_3', 'truncated', 's', '2018-09-12 11:20:20'),
(52, 'supplementary_3', 'In_reply_to_status_id_str', 's', '2018-09-12 11:20:20'),
(53, 'supplementary_3', 'In_reply_to_user_id', 's', '2018-09-12 11:20:20'),
(54, 'supplementary_3', 'In_reply_to_user_id_str', 's', '2018-09-12 11:20:20'),
(55, 'supplementary_3', 'In_reply_to_screen_name', 's', '2018-09-12 11:20:20'),
(56, 'supplementary_3', 'geo', 's', '2018-09-12 11:20:20'),
(57, 'supplementary_3', 'coordinates', 's', '2018-09-12 11:20:20'),
(58, 'supplementary_3', 'lang', 's', '2018-09-12 11:20:20'),
(59, 'supplementary_3', 'user_mentions', 's', '2018-09-12 11:20:20'),
(60, 'supplementary_3', 'place', 's', '2018-09-12 11:20:20'),
(61, 'supplementary_3', 'contributors', 's', '2018-09-12 11:20:20'),
(62, 'supplementary_3', 'is_quote_status', 's', '2018-09-12 11:20:20'),
(63, 'supplementary_3', 'quote_count', 's', '2018-09-12 11:20:20'),
(64, 'supplementary_3', 'favorited', 's', '2018-09-12 11:20:20'),
(65, 'supplementary_3', 'retweet_count', 's', '2018-09-12 11:20:20'),
(66, 'supplementary_3', 'retweeted', 's', '2018-09-12 11:20:20'),
(67, 'supplementary_3', 'filter_level', 's', '2018-09-12 11:20:20'),
(68, 'supplementary_3', 'timestamp_ms', 's', '2018-09-12 11:20:20'),
(69, 'supplementary_3', 'hashtags', 's', '2018-09-12 11:20:20'),
(70, 'supplementary_3', 'symbols', 's', '2018-09-12 11:20:20'),
(71, 'supplementary_3', 'timestamp', 's', '2018-09-12 11:20:20'),
(72, 'supplementary_4', 'PK', 's', '2018-09-12 11:20:45'),
(73, 'supplementary_4', 'app_id', 's', '2018-09-12 11:20:45'),
(74, 'supplementary_4', 'row_id', 's', '2018-09-12 11:20:45'),
(75, 'supplementary_4', 'id_str', 's', '2018-09-12 11:20:45'),
(76, 'supplementary_4', 'name', 's', '2018-09-12 11:20:45'),
(77, 'supplementary_4', 'location', 's', '2018-09-12 11:20:45'),
(78, 'supplementary_4', 'url', 's', '2018-09-12 11:20:45'),
(79, 'supplementary_4', 'description', 's', '2018-09-12 11:20:45'),
(80, 'supplementary_4', 'translatior_type', 's', '2018-09-12 11:20:45'),
(81, 'supplementary_4', 'protected', 's', '2018-09-12 11:20:45'),
(82, 'supplementary_4', 'verified', 's', '2018-09-12 11:20:45'),
(83, 'supplementary_4', 'followers_count', 's', '2018-09-12 11:20:45'),
(84, 'supplementary_4', 'friends_count', 's', '2018-09-12 11:20:45'),
(85, 'supplementary_4', 'listed_count', 's', '2018-09-12 11:20:45'),
(86, 'supplementary_4', 'favourites_count', 's', '2018-09-12 11:20:45'),
(87, 'supplementary_4', 'statuses_count', 's', '2018-09-12 11:20:45'),
(88, 'supplementary_4', 'utc_offset', 's', '2018-09-12 11:20:45'),
(89, 'supplementary_4', 'time_zone', 's', '2018-09-12 11:20:45'),
(90, 'supplementary_4', 'geo_enabled', 's', '2018-09-12 11:20:45'),
(91, 'supplementary_4', 'lang', 's', '2018-09-12 11:20:45'),
(92, 'supplementary_4', 'contributors_enabled', 's', '2018-09-12 11:20:45'),
(93, 'supplementary_4', 'is_translator', 's', '2018-09-12 11:20:45'),
(94, 'supplementary_4', 'profile_background_image_url_https', 's', '2018-09-12 11:20:45'),
(95, 'supplementary_4', ' profile_background_title', 's', '2018-09-12 11:20:45'),
(96, 'supplementary_4', 'profile_link_color', 's', '2018-09-12 11:20:45'),
(97, 'supplementary_4', 'profile_sidebar_border_color', 's', '2018-09-12 11:20:45'),
(98, 'supplementary_4', 'profile_sidebar_fill_color', 's', '2018-09-12 11:20:45'),
(99, 'supplementary_4', 'profile_text_color', 's', '2018-09-12 11:20:45'),
(100, 'supplementary_4', 'profile_use_background_image', 's', '2018-09-12 11:20:45'),
(101, 'supplementary_4', 'profile_image_url', 's', '2018-09-12 11:20:45'),
(102, 'supplementary_4', 'profile_image_url_https', 's', '2018-09-12 11:20:45'),
(103, 'supplementary_4', 'default_profile', 's', '2018-09-12 11:20:45'),
(104, 'supplementary_4', 'default_profile_image', 's', '2018-09-12 11:20:45'),
(105, 'supplementary_4', 'following', 's', '2018-09-12 11:20:45'),
(106, 'supplementary_4', 'following_request_sent', 's', '2018-09-12 11:20:45'),
(107, 'supplementary_4', 'notification', 's', '2018-09-12 11:20:45'),
(108, 'supplementary_4', 'timestamp', 's', '2018-09-12 11:20:45'),
(109, 'supplementary_5', 'PK', 's', '2018-09-12 11:21:12'),
(110, 'supplementary_5', 'app_id', 's', '2018-09-12 11:21:12'),
(111, 'supplementary_5', 'row_id', 's', '2018-09-12 11:21:12'),
(112, 'supplementary_5', 'link_karma', 's', '2018-09-12 11:21:12'),
(113, 'supplementary_5', 'Retrieved_On', 's', '2018-09-12 11:21:12'),
(114, 'supplementary_5', 'comment_karma', 's', '2018-09-12 11:21:12'),
(115, 'supplementary_5', 'profile_over_18', 's', '2018-09-12 11:21:12'),
(116, 'supplementary_5', 'timestamp', 's', '2018-09-12 11:21:12'),
(117, 'supplementary_6', 'PK', 's', '2018-09-12 11:21:29'),
(118, 'supplementary_6', 'app_id', 's', '2018-09-12 11:21:29'),
(119, 'supplementary_6', 'row_id', 's', '2018-09-12 11:21:29'),
(120, 'supplementary_6', 'Type', 's', '2018-09-12 11:21:29'),
(121, 'supplementary_6', 'author_url', 's', '2018-09-12 11:21:29'),
(122, 'supplementary_6', 'thumbernail_height', 's', '2018-09-12 11:21:29'),
(123, 'supplementary_6', 'thumbernail_url', 's', '2018-09-12 11:21:29'),
(124, 'supplementary_6', 'html', 's', '2018-09-12 11:21:29'),
(125, 'supplementary_6', 'Author_name', 's', '2018-09-12 11:21:29'),
(126, 'supplementary_6', 'provider_name', 's', '2018-09-12 11:21:29'),
(127, 'supplementary_6', 'Title', 's', '2018-09-12 11:21:29'),
(128, 'supplementary_6', 'provider_url', 's', '2018-09-12 11:21:29'),
(129, 'supplementary_6', 'version', 's', '2018-09-12 11:21:29'),
(130, 'supplementary_6', 'thumbnail_width', 's', '2018-09-12 11:21:29'),
(131, 'supplementary_6', 'width', 's', '2018-09-12 11:21:29'),
(132, 'supplementary_6', 'height', 's', '2018-09-12 11:21:29'),
(133, 'supplementary_6', 'timestamp', 's', '2018-09-12 11:21:29'),
(134, 'supplementary_7', 'PK', 's', '2018-09-12 11:21:49'),
(135, 'supplementary_7', 'app_id', 's', '2018-09-12 11:21:49'),
(136, 'supplementary_7', 'row_id', 's', '2018-09-12 11:21:49'),
(137, 'supplementary_7', 'height', 's', '2018-09-12 11:21:49'),
(138, 'supplementary_7', 'scrolling', 's', '2018-09-12 11:21:49'),
(139, 'supplementary_7', 'width', 's', '2018-09-12 11:21:49'),
(140, 'supplementary_7', 'content', 's', '2018-09-12 11:21:49'),
(141, 'supplementary_7', 'timestamp', 's', '2018-09-12 11:21:49'),
(142, 'supplementary_8', 'PK', 's', '2018-09-12 11:22:10'),
(143, 'supplementary_8', 'app_id', 's', '2018-09-12 11:22:10'),
(144, 'supplementary_8', 'row_id', 's', '2018-09-12 11:22:10'),
(145, 'supplementary_8', 'link_karma', 's', '2018-09-12 11:22:10'),
(146, 'supplementary_8', 'Image_id', 's', '2018-09-12 11:22:10'),
(147, 'supplementary_8', 'width', 's', '2018-09-12 11:22:10'),
(148, 'supplementary_8', 'height', 's', '2018-09-12 11:22:10'),
(149, 'supplementary_8', 'url', 's', '2018-09-12 11:22:10'),
(150, 'supplementary_8', 'timestamp', 's', '2018-09-12 11:22:10'),
(151, 'supplementary_9', 'PK', 's', '2018-09-12 11:22:44'),
(152, 'supplementary_9', 'app_id', 's', '2018-09-12 11:22:44'),
(153, 'supplementary_9', 'row_id', 's', '2018-09-12 11:22:44'),
(154, 'supplementary_9', 'Image_id', 's', '2018-09-12 11:22:44'),
(155, 'supplementary_9', 'Source_width', 's', '2018-09-12 11:22:44'),
(156, 'supplementary_9', 'Source_height', 's', '2018-09-12 11:22:44'),
(157, 'supplementary_9', 'Source_url', 's', '2018-09-12 11:22:44'),
(158, 'supplementary_9', 'Variants', 's', '2018-09-12 11:22:44'),
(159, 'supplementary_9', 'timestamp', 's', '2018-09-12 11:22:44'),
(160, 'supplementary_10', 'PK', 's', '2018-09-12 11:22:59'),
(161, 'supplementary_10', 'app_id', 's', '2018-09-12 11:22:59'),
(162, 'supplementary_10', 'row_id', 's', '2018-09-12 11:22:59'),
(163, 'supplementary_10', 'archived', 's', '2018-09-12 11:22:59'),
(164, 'supplementary_10', 'link_flair_text', 's', '2018-09-12 11:22:59'),
(165, 'supplementary_10', 'Saved', 's', '2018-09-12 11:22:59'),
(166, 'supplementary_10', 'Thumb_nail', 's', '2018-09-12 11:22:59'),
(167, 'supplementary_10', 'link_flair_css_class', 's', '2018-09-12 11:22:59'),
(168, 'supplementary_10', 'spoiler', 's', '2018-09-12 11:22:59'),
(169, 'supplementary_10', 'edited', 's', '2018-09-12 11:22:59'),
(170, 'supplementary_10', 'domain', 's', '2018-09-12 11:22:59'),
(171, 'supplementary_10', 'hide_score', 's', '2018-09-12 11:22:59'),
(172, 'supplementary_10', 'contest_mode', 's', '2018-09-12 11:22:59'),
(173, 'supplementary_10', 'permalink', 's', '2018-09-12 11:22:59'),
(174, 'supplementary_10', 'distinguished', 's', '2018-09-12 11:22:59'),
(175, 'supplementary_10', 'subreddit_id', 's', '2018-09-12 11:22:59'),
(176, 'supplementary_10', 'name', 's', '2018-09-12 11:22:59'),
(177, 'supplementary_10', 'locked', 's', '2018-09-12 11:22:59'),
(178, 'supplementary_10', 'gilded', 's', '2018-09-12 11:22:59'),
(179, 'supplementary_10', 'subreddit', 's', '2018-09-12 11:22:59'),
(180, 'supplementary_10', 'over_18', 's', '2018-09-12 11:22:59'),
(181, 'supplementary_10', 'media_embed', 's', '2018-09-12 11:22:59'),
(182, 'supplementary_10', 'is_itself', 's', '2018-09-12 11:22:59'),
(183, 'supplementary_10', ' author_flair_text', 's', '2018-09-12 11:22:59'),
(184, 'supplementary_10', 'stickied', 's', '2018-09-12 11:22:59'),
(185, 'supplementary_10', 'num_comments', 's', '2018-09-12 11:22:59'),
(186, 'supplementary_10', 'secure_media_embed', 's', '2018-09-12 11:22:59'),
(187, 'supplementary_10', 'quarantine', 's', '2018-09-12 11:22:59'),
(188, 'supplementary_10', 'preview', 's', '2018-09-12 11:22:59'),
(189, 'supplementary_10', 'post_hint', 's', '2018-09-12 11:22:59'),
(190, 'supplementary_10', 'Downs', 's', '2018-09-12 11:22:59'),
(191, 'supplementary_10', 'author_flair_css_class', 's', '2018-09-12 11:22:59'),
(192, 'supplementary_10', 'timestamp', 's', '2018-09-12 11:22:59');

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
  ADD KEY `FK_apps` (`app_id`),
  ADD KEY `FK_data_type` (`data_type`);

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
-- Indexes for table `physical_schemas`
--
ALTER TABLE `physical_schemas`
  ADD PRIMARY KEY (`PK`);

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
-- AUTO_INCREMENT for table `physical_schemas`
--
ALTER TABLE `physical_schemas`
  MODIFY `PK` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=193;

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
-- Constraints for dumped tables
--

--
-- Constraints for table `app_schemas`
--
ALTER TABLE `app_schemas`
  ADD CONSTRAINT `FK_apps` FOREIGN KEY (`app_id`) REFERENCES `apps` (`PK`),
  ADD CONSTRAINT `FK_data_type` FOREIGN KEY (`data_type`) REFERENCES `data_types` (`PK`);

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
  ADD CONSTRAINT `FK_physical_attr` FOREIGN KEY (`physical_attribute`) REFERENCES `physical_schemas` (`PK`);

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
