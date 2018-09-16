import MySQLdb

db_name = "stencil_storage"
table_name = "supplementary_10"

db_conn = MySQLdb.connect(
        host   = "127.0.0.1",
        port   = 3307,
        user   = "root",
        passwd = "",
        db     = db_name,
    )

cur = db_conn.cursor()

q = "SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA`='%s' AND `TABLE_NAME`='%s';" % (db_name, table_name)

cur.execute(q)
tables = cur.fetchall()

q = "INSERT INTO `physical_schemas` (`table_name`, `column_name`, `type`) VALUES "

for table in tables:
    q += "( '%s', '%s', 's' ), " % (table_name, table[0])

print q[:-2]