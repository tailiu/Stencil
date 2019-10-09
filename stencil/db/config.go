package db

import "database/sql"

var dbConns map[string]*sql.DB

const STENCIL_DB = "stencil"
const DB_ADDR = "10.230.12.86"
const DB_ADDR_old = "10.230.12.75"
const DB_PORT = "5432"
const DB_USER = "cow"
const DB_PASSWORD = "123456"
const OTHER_SERVER_ADDR = "10.230.12.75"
const OTHER_SERVER_PORT = "4410"
const OTHER_SERVER_MEDIA_PATH = "/home/zain/project/resources/"