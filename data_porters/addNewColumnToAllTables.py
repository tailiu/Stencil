import psycopg2
import threading 

def getDBConn(db):
    conn = psycopg2.connect(dbname=db, user="cow", password="123456", host="10.230.12.86", port="5432")
    cursor = conn.cursor()
    return conn, cursor

def getTables():
    conn, cur = getDBConn(dbName)
    tables = []
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    cur.execute(tableq)
    for row in cur.fetchall():
        table = row[0]
        tables.append(table)
    return tables

def execQuery(tableName):
    conn, cur = getDBConn(dbName)
    with conn:
        q = "ALTER TABLE \"%s\" ADD display_flag bool DEFAULT false;" %tableName
        print q
        cur.execute(q)
        conn.commit()
        print "===>> %s finished!!!" % tableName

if __name__ == "__main__":
    
    dbName  = "diaspora_1000"
    threads = []

    for table in getTables():
        threads.append(threading.Thread(target=execQuery, args=[table]))
    
    for x in threads: x.start()
    for x in threads: x.join()
        