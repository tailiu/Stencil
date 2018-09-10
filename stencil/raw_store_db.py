import MySQLdb, json

def getDBConn():
    db_conn = MySQLdb.connect(
        host   = "127.0.0.1",
        port   = 3307,
        user   = "root",
        passwd = "",
        db     = "stencil_storage",
    )
    return db_conn, db_conn.cursor()

def getSchemaMapping(app_name, table_name):
    sql =  "SELECT GROUP_CONCAT(LOWER(app_mappings.column_name), ',', LOWER(app_mappings.mapping)) \
            FROM app_mappings JOIN apps \
            ON app_mappings.app_id = apps.PK \
            JOIN  app_tables \
            ON app_mappings.table_id = app_tables.PK \
            WHERE apps.app_name = '%s' AND app_tables.table_name = '%s'" % (app_name, table_name)
    CUR.execute(sql)
    result = CUR.fetchone()[0].split(",")
    return [(result[i], result[i+1]) for i in range(0, len(result), 2)]

################ DB Globals ##
CONN, CUR = getDBConn()
##############################

if __name__ == "__main__":

    hn_fpath = "/Users/zain/Documents/DataSets/HackerNews/hn.min2.json"

    with open(hn_fpath) as fh: data = json.load(fh)
    
    sql = "INSERT INTO Story "

    hn_story_schema = getSchemaMapping("hacker news", "story")
    hn_comment_schema = getSchemaMapping("hacker news", "comment")

    for datum in data:
        if datum['type'] == "story":
            sql += "( "
            attrs = [ x.lower() for x in datum.keys() if x.lower() in [y[0] for y in hn_story_schema] ]
            for attr in attrs:
                sql += str(attr) + ", "
            sql = sql[:-2] + " ) VALUES ( "
            for attr in attrs:
                sql += '"' + str(datum[attr]) + '", '
            sql = sql[:-2] + ")"
            print sql
            break

    pass