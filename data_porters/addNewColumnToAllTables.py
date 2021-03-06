import psycopg2
import threading 

def getDBConn(db):
    conn = psycopg2.connect(dbname=db, user="cow", password="123456", host="10.230.12.86", port="5432")
    cursor = conn.cursor()
    return conn, cursor

def getTables():
    tables = []
    conn, cur = getDBConn(dbName)
    with conn:
        tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
        cur.execute(tableq)
        for row in cur.fetchall():
            table = row[0]
            tables.append(table)
    return tables

def execQuery(tableName):
    conn, cur = getDBConn(dbName)
    with conn:
        # q = "ALTER TABLE \"%s\" ADD display_flag bool DEFAULT false;" %tableName
        q = "ALTER TABLE \"%s\" ADD COLUMN IF NOT EXISTS id int8;" %tableName
        print ">> %s STARTED!!!" % tableName
        cur.execute(q)
        conn.commit()
        print "==>> %s FINISHED!!!" % tableName

if __name__ == "__main__":
    
    dbName  = "gnusocial_template"
    threads = []

    for table in getTables():
        threads.append(threading.Thread(target=execQuery, args=[table]))
    
    for x in threads: x.start()
    for x in threads: x.join()
        