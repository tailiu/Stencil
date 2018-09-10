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

def getSchema(app_name, table_name):
    sql = "SELECT GROUP_CONCAT(LOWER(app_schemas.column_name)) \
           FROM app_schemas JOIN apps \
           ON app_schemas.app_id = apps.PK \
           WHERE apps.app_name = '%s' AND app_schemas.table_name = '%s'" % (app_name, table_name)

    CUR.execute(sql)
    result = CUR.fetchone()[0].split(',')
    return result

def getMapping(app_name, table_name):
    sql = "SELECT GROUP_CONCAT(LOWER(app_schemas.column_name), LOWER(app_schemas.mapping)) \
           FROM app_schemas JOIN apps \
           ON app_schemas.app_id = apps.PK \
           WHERE apps.app_name = '%s' AND app_schemas.table_name = '%s'" % (app_name, table_name)

    CUR.execute(sql)
    result = CUR.fetchone()[0].split(',')
    return result

############### Globals #
CONN, CUR = getDBConn()
#########################

if __name__ == "__main__":

    # getMapping("hacker news", "story")
    # exit(0)

    hn_fpath = "/Users/zain/Documents/DataSets/HackerNews/hn.min2.json"

    with open(hn_fpath) as fh: data = json.load(fh)
    
    sql = "INSERT INTO Story "

    hn_story_schema = getSchema("hacker news", "story")

    for datum in data:
        if datum['type'] == "story":
            sql += "( "
            attrs = [ x.lower() for x in datum.keys() if x.lower() in hn_story_schema ]
            for attr in attrs:
                sql += str(attr) + ", "
            sql = sql[:-2] + " ) VALUES ( "
            for attr in attrs:
                sql += '"' + str(datum[attr]) + '", '
            sql = sql[:-2] + ")"
            print sql
            break

    pass