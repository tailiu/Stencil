import psycopg2
import apps
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT

def getDBConn(db):
    conn = psycopg2.connect(dbname=db, user="root", host="10.224.45.158", port="26257")
    conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT) 
    return conn, conn.cursor()

def createTables(tno):
    base_schema = """
            CREATE TABLE IF NOT EXISTS %s (
            row_id serial PRIMARY KEY,
            %s
        )""" % ("base_%s"%tno)

if __name__ == "__main__":

    for app in apps.getApps():
        db = getDBConn(app)
        schemas[app] = getTables(db.cur)
        db.cur.close()
        db.conn.close()