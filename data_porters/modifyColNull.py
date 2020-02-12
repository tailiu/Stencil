import psycopg2
import time

host = "10.230.12.86"

def diconnectOtherConns(dbName):
    conn = psycopg2.connect(dbname="stencil", user="cow", password="123456", host=host, port="5432")
    with conn:
        cursor = conn.cursor()
        query = "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();"%dbName
        cursor.execute(query)

def getDBConn(dbName):
    conn = psycopg2.connect(dbname=dbName, user="cow", password="123456", host=host, port="5432")
    cursor = conn.cursor()
    return conn, cursor

def getTables(dbName):
    tables = []
    conn, cur = getDBConn(dbName)
    with conn:
        tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
        cur.execute(tableq)
        for row in cur.fetchall():
            table = row[0]
            tables.append(table)
    return set(tables)

def execQuery(dbName, tableName, columns):
    print ">> %s STARTED!!!" % tableName
    conn, cur = getDBConn(dbName)
    with conn:
        for column in columns:
            q = "ALTER TABLE \"%s\" ALTER COLUMN \"%s\" DROP NOT NULL;"%(tableName, column)
            print q
            cur.execute(q)
        conn.commit()
    print ">> %s FINISHED!!!" % tableName

if __name__ == "__main__":
    
    for db in ["diaspora_test", "diaspora_10000", "diaspora_100000", "diaspora_1000000"]:
        diconnectOtherConns(db)
        tables = ["likes"]
        columns = ["created_at", "updated_at"]
        for table in tables:
            execQuery(db, table, columns)
        