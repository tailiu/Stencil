import MySQLdb, datetime

db_name = "stencil_storage"

db_conn = MySQLdb.connect(
        host   = "127.0.0.1",
        port   = 3307,
        user   = "root",
        passwd = "",
        db     = db_name,
    )

cur = db_conn.cursor()

# table_name = "supplementary_10"

# q = "SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA`='%s' AND `TABLE_NAME`='%s';" % (db_name, table_name)

# cur.execute(q)
# tables = cur.fetchall()

# q = "INSERT INTO `physical_schemas` (`table_name`, `column_name`, `type`) VALUES "

# for table in tables:
#     q += "( '%s', '%s', 's' ), " % (table_name, table[0])

# print q[:-2]


# fname = "./datasets/twitter.tweet.cols"

# with open(fname, "rb") as fh:
#     rows = fh.read().split("\n")
#     q   = "INSERT INTO `app_schemas` (`table_id`, `column_name`, `data_type`) VALUES ('4', '%s', '1');"
#     q   = "Create Table tweet ( `PK` int(11) NOT NULL, %s `Timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP )"
#     tables = ""
#     for row in rows:
#         tables += row.strip(' ",') +  " varchar(256) COLLATE utf8_unicode_ci DEFAULT NULL, "
#         # print q % row.strip(' ",')
#     print q % tables

pre_time = datetime.datetime.now().time()

sql = """ SELECT * 
          FROM base_1 JOIN 
               supplementary_3 ON base_1.row_id = supplementary_3.row_id JOIN
               app_tables ON app_tables.app_id = base_1.app_id
          WHERE base_1.app_id = 2 AND app_tables.table_name = "Tweet"
      """

cur.execute(sql)
post_time = datetime.datetime.now().time()

print post_time, pre_time, post_time - pre_time