package db

import "database/sql"

var dbConns map[string]*sql.DB

var STENCIL_DB = "stencil_exp2"
var DIASPORA_DB = "diaspora_1000000_exp2"
var MASTODON_DB = "mastodon_exp2"
const DB_TEST = false
const DB_ADDR = "10.230.12.86"
const DB_ADDR_old = "10.230.12.75"
const DB_PORT = "5432"
const DB_USER = "cow"
const DB_PASSWORD = "123456"
const FTP_USER = "cowftp"
const FTP_PASSWORD = "Big1Fat2Cow3"
const FTP_SERVER_ADDR = "10.230.12.75"
const FTP_SERVER_PORT = "21"
const FTP_SERVER_MEDIA_PATH = "/home/zain/project/resources/"
