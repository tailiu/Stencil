import psycopg2, copy
from psycopg2.extras import RealDictCursor

def getDBConn(db, cursor_dict=False):
    conn = psycopg2.connect(dbname=db, user="cow", password="123456", host="10.230.12.86", port="5432")
    if cursor_dict is True:
        cursor = conn.cursor(cursor_factory = RealDictCursor)
    else:
        cursor = conn.cursor()
    return conn, cursor


def getAppTables(app):
    tables = []
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    cur.execute(tableq)
    for row in cur.fetchall():
        tables.append(row["tablename"])
    return tables

def createIndices(app):
    column_name = "mark_as_delete"
    for table in getAppTables(app):
        idxsql = "CREATE INDEX %s_id_%s_idx ON public.%s (%s, id); "%(table, column_name, table, column_name)
        print idxsql
        try:
            cur.execute(idxsql)
            dbConn.commit()
        except Exception as e:
            print e
            dbConn.rollback()

if __name__ == "__main__":
    app = "diaspora"
    dbConn, cur = getDBConn(app, True)
    createIndices(app)