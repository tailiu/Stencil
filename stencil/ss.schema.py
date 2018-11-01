import psycopg2

# Connect to the "bank" database.
conn = psycopg2.connect(
    database='twitter',
    user='root',
    sslmode='disable',
    port=26257,
    host='10.224.45.162'
)

# Make each statement commit immediately.
conn.set_session(autocommit=True)

# Open a cursor to perform database operations.
cur = conn.cursor()

# sql = "DROP DATABASE stencil_storage CASCADE; CREATE DATABASE stencil_storage;"

sql = """
    CREATE TABLE apps (
  PK  serial PRIMARY KEY,
  app_name text  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

--
-- Dumping data for table apps
--

INSERT INTO apps (PK, app_name, timestamp) VALUES
(1, 'reddit', '2018-09-09 11:18:19'),
(2, 'twitter', '2018-09-09 11:18:19'),
(3, 'hacker news', '2018-09-09 11:18:19');

-- --------------------------------------------------------

--
-- Table structure for table app_schemas
--

CREATE TABLE app_schemas (
  PK  serial PRIMARY KEY,
  table_id  int NOT NULL,
  column_name text  NOT NULL,
  data_type  int NOT NULL,
  constraints text  DEFAULT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

--
-- Dumping data for table app_schemas
--

INSERT INTO app_schemas (PK, table_id, column_name, data_type, constraints, timestamp) VALUES
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
(21, 3, 'created_at', 1, NULL, '2018-09-19 11:47:26'),
(22, 3, 'id', 1, NULL, '2018-09-19 11:58:36'),
(23, 3, 'id_str', 1, NULL, '2018-09-19 11:58:36'),
(24, 3, 'text', 1, NULL, '2018-09-19 11:58:36'),
(25, 3, 'truncated', 1, NULL, '2018-09-19 11:58:36'),
(26, 3, 'in_reply_to_status_id', 1, NULL, '2018-09-19 11:58:36'),
(27, 3, 'in_reply_to_status_id_str', 1, NULL, '2018-09-19 11:58:36'),
(28, 3, 'in_reply_to_user_id', 1, NULL, '2018-09-19 11:58:36'),
(29, 3, 'in_reply_to_user_id_str', 1, NULL, '2018-09-19 11:58:36'),
(30, 3, 'in_reply_to_screen_name', 1, NULL, '2018-09-19 11:58:36'),
(31, 3, 'geo', 1, NULL, '2018-09-19 11:58:36'),
(32, 3, 'coordinates', 1, NULL, '2018-09-19 11:58:36'),
(33, 3, 'place', 1, NULL, '2018-09-19 11:58:36'),
(34, 3, 'contributors', 1, NULL, '2018-09-19 11:58:36'),
(35, 3, 'is_quote_status', 1, NULL, '2018-09-19 11:58:36'),
(36, 3, 'quote_count', 1, NULL, '2018-09-19 11:58:36'),
(37, 3, 'reply_count', 1, NULL, '2018-09-19 11:58:36'),
(38, 3, 'retweet_count', 1, NULL, '2018-09-19 11:58:36'),
(39, 3, 'favorite_count', 1, NULL, '2018-09-19 11:58:36'),
(40, 3, 'favorited', 1, NULL, '2018-09-19 11:58:36'),
(41, 3, 'retweeted', 1, NULL, '2018-09-19 11:58:36'),
(42, 3, 'filter_level', 1, NULL, '2018-09-19 11:58:36'),
(43, 3, 'lang', 1, NULL, '2018-09-19 11:58:36'),
(44, 3, 'Timestamp_ms', 1, NULL, '2018-09-19 11:58:36'),
(45, 3, 'hashtags', 1, NULL, '2018-09-19 11:58:36'),
(46, 3, 'urls', 1, NULL, '2018-09-19 11:58:36'),
(47, 3, 'user_mentions', 1, NULL, '2018-09-19 11:58:36'),
(48, 3, 'symbols', 1, NULL, '2018-09-19 11:58:36'),
(49, 4, 'Id', 1, NULL, '2018-09-19 11:59:46'),
(50, 4, 'id_str', 1, NULL, '2018-09-19 11:59:46'),
(51, 4, 'name', 1, NULL, '2018-09-19 11:59:46'),
(52, 4, 'screen_name', 1, NULL, '2018-09-19 11:59:46'),
(53, 4, 'location', 1, NULL, '2018-09-19 11:59:46'),
(54, 4, 'url', 1, NULL, '2018-09-19 11:59:46'),
(55, 4, 'description', 1, NULL, '2018-09-19 11:59:46'),
(56, 4, 'translator_type', 1, NULL, '2018-09-19 11:59:46'),
(57, 4, 'protected', 1, NULL, '2018-09-19 11:59:46'),
(58, 4, 'verified', 1, NULL, '2018-09-19 11:59:46'),
(59, 4, 'followers_count', 1, NULL, '2018-09-19 11:59:46'),
(60, 4, 'friends_count', 1, NULL, '2018-09-19 11:59:46'),
(61, 4, 'listed_count', 1, NULL, '2018-09-19 11:59:46'),
(62, 4, 'favourites_count', 1, NULL, '2018-09-19 11:59:46'),
(63, 4, 'statuses_count', 1, NULL, '2018-09-19 11:59:46'),
(64, 4, 'created_at', 1, NULL, '2018-09-19 11:59:46'),
(65, 4, 'utc_offset', 1, NULL, '2018-09-19 11:59:46'),
(66, 4, 'time_zone', 1, NULL, '2018-09-19 11:59:46'),
(67, 4, 'geo_enabled', 1, NULL, '2018-09-19 11:59:46'),
(68, 4, 'lang', 1, NULL, '2018-09-19 11:59:46'),
(69, 4, 'contributors_enabled', 1, NULL, '2018-09-19 11:59:46'),
(70, 4, 'is_translator', 1, NULL, '2018-09-19 11:59:46'),
(71, 4, 'profile_background_color', 1, NULL, '2018-09-19 11:59:46'),
(72, 4, 'profile_background_image_url', 1, NULL, '2018-09-19 11:59:46'),
(73, 4, 'profile_background_image_url_https', 1, NULL, '2018-09-19 11:59:46'),
(74, 4, 'profile_background_tile', 1, NULL, '2018-09-19 11:59:46'),
(75, 4, 'profile_link_color', 1, NULL, '2018-09-19 11:59:46'),
(76, 4, 'profile_sidebar_border_color', 1, NULL, '2018-09-19 11:59:46'),
(77, 4, 'profile_sidebar_fill_color', 1, NULL, '2018-09-19 11:59:46'),
(78, 4, 'profile_text_color', 1, NULL, '2018-09-19 11:59:46'),
(79, 4, 'profile_use_background_image', 1, NULL, '2018-09-19 11:59:46'),
(80, 4, 'profile_image_url', 1, NULL, '2018-09-19 11:59:46'),
(81, 4, 'profile_image_url_https', 1, NULL, '2018-09-19 11:59:46'),
(82, 4, 'default_profile', 1, NULL, '2018-09-19 11:59:46'),
(83, 4, 'default_profile_image', 1, NULL, '2018-09-19 11:59:46'),
(84, 4, 'following', 1, NULL, '2018-09-19 11:59:46'),
(85, 4, 'follow_request_sent', 1, NULL, '2018-09-19 11:59:46'),
(86, 4, 'notifications', 1, NULL, '2018-09-19 11:59:46'),
(87, 3, 'user', 1, NULL, '2018-09-19 12:22:05');

-- --------------------------------------------------------

--
-- Table structure for table app_tables
--

CREATE TABLE app_tables (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  table_name text  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

--
-- Dumping data for table app_tables
--

INSERT INTO app_tables (PK, app_id, table_name, timestamp) VALUES
(1, 3, 'Story', '2018-09-16 10:09:44'),
(2, 3, 'Comment', '2018-09-16 10:09:44'),
(3, 2, 'Tweet', '2018-09-19 11:46:23'),
(4, 2, 'User', '2018-09-19 11:46:23');

-- --------------------------------------------------------

--
-- Table structure for table base_1
--

CREATE TABLE base_1 (
  PK  serial PRIMARY KEY,
  app_id  int DEFAULT NULL,
  row_id varchar(128)  DEFAULT NULL,
  "User" text ,
  Time text ,
  URL text ,
  Parent text ,
  Text text ,
  Id text ,
  Retrieved_On text ,
  Score text ,
  Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table base_2
--

CREATE TABLE base_2 (
  PK  serial PRIMARY KEY,
  app_id  int DEFAULT NULL,
  row_id varchar(128)  DEFAULT NULL,
  Ups text ,
  Num_Comments text ,
  Name text ,
  Profile_IMG text ,
  Profile_Color text ,
  Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table base_table_attributes
--

CREATE TABLE base_table_attributes (
  PK  serial PRIMARY KEY,
  table_name text  NOT NULL,
  column_name text  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

--
-- Dumping data for table base_table_attributes
--

INSERT INTO base_table_attributes (PK, table_name, column_name, timestamp) VALUES
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
-- Table structure for table data_types
--

CREATE TABLE data_types (
  PK  serial PRIMARY KEY,
  data_type varchar(128)  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

--
-- Dumping data for table data_types
--

INSERT INTO data_types (PK, data_type, timestamp) VALUES
(1, 'text', '2018-09-12 07:07:27'),
(2, 'int', '2018-09-12 07:07:27');

-- --------------------------------------------------------

--
-- Table structure for table physical_mappings
--

CREATE TABLE physical_mappings (
  PK  serial PRIMARY KEY,
  logical_attribute  int DEFAULT NULL,
  physical_attribute  int DEFAULT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

--
-- Dumping data for table physical_mappings
--

INSERT INTO physical_mappings (PK, logical_attribute, physical_attribute, timestamp) VALUES
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
(17, 18, 20, '2018-09-12 12:01:51'),
(19, 87, 16, '2018-09-19 12:23:19'),
(20, 21, 17, '2018-09-19 12:33:21'),
(21, 46, 18, '2018-09-19 12:33:21'),
(22, 26, 19, '2018-09-19 12:33:21'),
(23, 24, 20, '2018-09-19 12:33:21'),
(24, 22, 21, '2018-09-19 12:33:21'),
(25, 39, 28, '2018-09-19 12:33:21'),
(26, 37, 29, '2018-09-19 12:33:21'),
(27, 52, 30, '2018-09-19 12:33:21'),
(28, 72, 31, '2018-09-19 12:33:21'),
(29, 71, 32, '2018-09-19 12:33:21');

-- --------------------------------------------------------

--
-- Table structure for table schema_mappings
--

CREATE TABLE schema_mappings (
  PK  serial PRIMARY KEY,
  app1_attribute  int NOT NULL,
  app2_attribute  int NOT NULL,
  rules text  DEFAULT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_1
--

CREATE TABLE supplementary_1 (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  Descendents text ,
  Type text ,
  Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_2
--

CREATE TABLE supplementary_2 (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  Kids text ,
  Retrieved_On text ,
  Type text ,
  Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_3
--

CREATE TABLE supplementary_3 (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  id_str text ,
  truncated text ,
  In_reply_to_status_id_str text ,
  In_reply_to_user_id text ,
  In_reply_to_user_id_str text ,
  In_reply_to_screen_name text ,
  geo text ,
  coordinates text ,
  lang text ,
  user_mentions text ,
  place text ,
  contributors text ,
  is_quote_status text ,
  quote_count text ,
  favorited text ,
  retweet_count text ,
  retweeted text ,
  filter_level text ,
  timestamp_ms text ,
  hashtags text ,
  symbols text ,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_4
--

CREATE TABLE supplementary_4 (
  PK  serial PRIMARY KEY ,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  id_str text ,
  name text ,
  location text ,
  url text ,
  description text ,
  translator_type text ,
  protected text ,
  verified text ,
  followers_count text ,
  friends_count text ,
  listed_count text ,
  favourites_count text ,
  statuses_count text ,
  utc_offset text ,
  time_zone text ,
  geo_enabled text ,
  lang text ,
  contributors_enabled text ,
  is_translator text ,
  profile_background_image_url_https text ,
  profile_background_tile text ,
  profile_link_color text ,
  profile_sidebar_border_color text ,
  profile_sidebar_fill_color text ,
  profile_text_color text ,
  profile_use_background_image text ,
  profile_image_url text ,
  profile_image_url_https text ,
  default_profile text ,
  default_profile_image text ,
  following text ,
  follow_request_sent text ,
  notifications text ,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_5
--

CREATE TABLE supplementary_5 (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  link_karma text ,
  Retrieved_On text ,
  comment_karma text ,
  profile_over_18 text ,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_6
--

CREATE TABLE supplementary_6 (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  Type text ,
  author_url text ,
  thumbernail_height text ,
  thumbernail_url text ,
  html text ,
  Author_name text ,
  provider_name text ,
  Title text ,
  provider_url text ,
  version text ,
  thumbnail_width text ,
  width text ,
  height text ,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_7
--

CREATE TABLE supplementary_7 (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  height text ,
  scrolling text ,
  width text ,
  content text ,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_8
--

CREATE TABLE supplementary_8 (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  link_karma text ,
  Image_id text ,
  width text ,
  height text ,
  url text ,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_9
--

CREATE TABLE supplementary_9 (
  PK  serial PRIMARY KEY,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  Image_id text ,
  Source_width text ,
  Source_height text ,
  Source_url text ,
  Variants text ,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_10
--

CREATE TABLE supplementary_10 (
  PK  serial PRIMARY KEY ,
  app_id  int NOT NULL,
  row_id varchar(128)  DEFAULT NULL,
  archived text ,
  link_flair_text text ,
  Saved text ,
  Thumb_nail text ,
  link_flair_css_class text ,
  spoiler text ,
  edited text ,
  domain text ,
  hide_score text ,
  contest_mode text ,
  permalink text ,
  distinguished text ,
  subreddit_id text ,
  name text ,
  locked text ,
  gilded text ,
  subreddit text ,
  over_18 text ,
  media_embed text ,
  is_itself text ,
   author_flair_text text ,
  stickied text ,
  num_comments text ,
  secure_media_embed text ,
  quarantine text ,
  preview text ,
  post_hint text ,
  Downs text ,
  author_flair_css_class text ,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)  ;

-- --------------------------------------------------------

--
-- Table structure for table supplementary_tables
--

CREATE TABLE supplementary_tables (
  PK  serial PRIMARY KEY,
  table_id  int NOT NULL,
  supplementary_table text  NOT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP 
)  ;

--
-- Dumping data for table supplementary_tables
--

INSERT INTO supplementary_tables (PK, table_id, supplementary_table, timestamp) VALUES
(1, 1, 'supplementary_1', '2018-09-16 10:23:19'),
(2, 2, 'supplementary_2', '2018-09-16 10:23:19'),
(3, 3, 'supplementary_3', '2018-09-19 12:01:33'),
(4, 4, 'supplementary_4', '2018-09-19 12:01:33');

"""

hn_schema = """
CREATE TABLE comment (
  PK serial PRIMARY KEY,
  By text DEFAULT NULL,
  Id text DEFAULT NULL,
  Retrieved_on varchar(128) DEFAULT NULL,
  Time varchar(128) DEFAULT NULL,
  Kids text DEFAULT NULL,
  Parent text DEFAULT NULL,
  Text text DEFAULT NULL,
  Type varchar(64) DEFAULT NULL,
  Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ;


CREATE TABLE story (
  PK serial PRIMARY KEY,
  By text DEFAULT NULL,
  Descendents text DEFAULT NULL,
  Id text DEFAULT NULL,
  Retrieved_on varchar(128) DEFAULT NULL,
  Score varchar(10) DEFAULT NULL,
  Time varchar(128) DEFAULT NULL,
  Title text,
  Type varchar(64) DEFAULT NULL,
  Url text DEFAULT NULL,
  Timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ;
"""

twitter_schema = """
  CREATE TABLE tweet (
  PK serial PRIMARY KEY,
  "user" text  DEFAULT NULL,
  created_at text  DEFAULT NULL,
  id text  DEFAULT NULL,
  id_str text  DEFAULT NULL,
  text text ,
  truncated text  DEFAULT NULL,
  in_reply_to_status_id text  DEFAULT NULL,
  in_reply_to_status_id_str text  DEFAULT NULL,
  in_reply_to_user_id text  DEFAULT NULL,
  in_reply_to_user_id_str text  DEFAULT NULL,
  in_reply_to_screen_name text  DEFAULT NULL,
  geo text  DEFAULT NULL,
  coordinates text  DEFAULT NULL,
  place text  DEFAULT NULL,
  contributors text ,
  is_quote_status text  DEFAULT NULL,
  quote_count text  DEFAULT NULL,
  reply_count text  DEFAULT NULL,
  retweet_count text  DEFAULT NULL,
  favorite_count text  DEFAULT NULL,
  favorited text  DEFAULT NULL,
  retweeted text  DEFAULT NULL,
  filter_level text  DEFAULT NULL,
  lang text  DEFAULT NULL,
  timestamp_ms text  DEFAULT NULL,
  hashtags text ,
  urls text  DEFAULT NULL,
  user_mentions text ,
  symbols text  DEFAULT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ;

-- --------------------------------------------------------

--
-- Table structure for table user
--

CREATE TABLE "user" (
  PK serial PRIMARY KEY,
  id text  DEFAULT NULL,
  id_str text  DEFAULT NULL,
  name text  DEFAULT NULL,
  screen_name text  DEFAULT NULL,
  location text  DEFAULT NULL,
  url text  DEFAULT NULL,
  description text ,
  translator_type text  DEFAULT NULL,
  protected text  DEFAULT NULL,
  verified text  DEFAULT NULL,
  followers_count text  DEFAULT NULL,
  friends_count text  DEFAULT NULL,
  listed_count text  DEFAULT NULL,
  favourites_count text  DEFAULT NULL,
  statuses_count text  DEFAULT NULL,
  created_at text  DEFAULT NULL,
  utc_offset text  DEFAULT NULL,
  time_zone text  DEFAULT NULL,
  geo_enabled text  DEFAULT NULL,
  lang text  DEFAULT NULL,
  contributors_enabled text  DEFAULT NULL,
  is_translator text  DEFAULT NULL,
  profile_background_color text  DEFAULT NULL,
  profile_background_image_url text  DEFAULT NULL,
  profile_background_image_url_https text  DEFAULT NULL,
  profile_background_tile text  DEFAULT NULL,
  profile_link_color text  DEFAULT NULL,
  profile_sidebar_border_color text  DEFAULT NULL,
  profile_sidebar_fill_color text  DEFAULT NULL,
  profile_text_color text  DEFAULT NULL,
  profile_use_background_image text  DEFAULT NULL,
  profile_image_url text  DEFAULT NULL,
  profile_image_url_https text  DEFAULT NULL,
  default_profile text  DEFAULT NULL,
  default_profile_image text  DEFAULT NULL,
  following text  DEFAULT NULL,
  follow_request_sent text  DEFAULT NULL,
  notifications text  DEFAULT NULL,
  timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ;

"""

# print sql.lower()
# exit()

cur.execute(twitter_schema.lower())

# # Create the "accounts" table.
# cur.execute("CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT)")

# # Insert two rows into the "accounts" table.
# cur.execute("INSERT INTO accounts (id, balance) VALUES (1, 1000), (2, 250)")

# # Print out the balances.
# cur.execute("SELECT id, balance FROM accounts")
# rows = cur.fetchall()
# print('Initial balances:')
# for row in rows:
#     print([str(cell) for cell in row])

# Close the database connection.
cur.close()
conn.close()