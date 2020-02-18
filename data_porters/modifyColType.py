import psycopg2
import threading 
import time

def diconnectOtherConns(dbName):
    conn = psycopg2.connect(dbname="stencil", user="cow", password="123456", host="10.230.12.86", port="5432")
    with conn:
        cursor = conn.cursor()
        query = "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();"%dbName
        cursor.execute(query)

def getDBConn(dbName):
    conn = psycopg2.connect(dbname=dbName, user="cow", password="123456", host="10.230.12.86", port="5432")
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

def getColumns(dbName, table):    
    conn, cur = getDBConn(dbName)
    columns = []
    with conn:
        q = "SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = 'public' AND table_name   = '%s';"%table
        cur.execute(q)
        for row in cur.fetchall():
            column_name = row[0]
            data_type = row[1]
            if data_type == "integer":
                columns.append(column_name)
    return columns

def execQuery(dbName, tableName):
    print ">> %s STARTED!!!" % tableName
    columns = getColumns(dbName, tableName)
    if len(columns) > 0:
        conn, cur = getDBConn(dbName)
        with conn:
            for column in columns:
                q = "ALTER TABLE \"%s\" ALTER COLUMN \"%s\" TYPE int8;" %(tableName, column)
                print q
                cur.execute(q)
            conn.commit()
            # conn.close()
    print ">> %s FINISHED!!!" % tableName

if __name__ == "__main__":
    
    for db in ["gnusocial_test", "gnusocial_template", "twitter_test", "twitter_template"]:
        print "\n\nCURRENT DB : %s\n\n"%db
        time.sleep(3)
        diconnectOtherConns(db)
        threads = []
        for table in getTables(db):
            execQuery(db, table)
        #     threads.append(threading.Thread(target=execQuery, args=[db,table]))
        # for x in threads: x.start()
        # for x in threads: x.join()
        # del threads[:]
        # break
        