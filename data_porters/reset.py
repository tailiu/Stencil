import psycopg2, datetime, time

def getDB(dbname, blade=False):
    if blade:
        host = "10.230.12.75"
    else:
        host = "10.230.12.86"
    conn = psycopg2.connect(dbname=dbname, user="cow", host=host, port="5432", password="123456")
    cursor = conn.cursor()
    return conn, cursor

def truncate(dbname, blade):
    conn, cur = getDB(dbname, blade)
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    cur.execute(tableq)
    for row in cur.fetchall():
        q = 'TRUNCATE "%s" RESTART IDENTITY CASCADE;'%row[0]
        print q
        cur.execute(q)
    conn.commit()

def reverseMarkAsDelete(dbname, blade):
    conn, cur = getDB(dbname, blade)
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    cur.execute(tableq)
    for row in cur.fetchall():
        q = 'UPDATE "%s" SET mark_as_delete = false ;'%row[0]
        print q
        cur.execute(q)
    conn.commit()

def truncateTableFromStencil(table):
    conn, cur = getDB("stencil", blade=False)
    q = 'TRUNCATE "%s" RESTART IDENTITY CASCADE;'%table
    print q
    cur.execute(q)
    conn.commit()

def resetRowDesc():
    conn, cur = getDB("stencil", blade=False)
    
    q = "select rowid from row_desc group by rowid having count(*) > 1;"
    cur.execute(q)
    for row in cur.fetchall():
        q = "delete from row_desc where app_id != 1 and rowid = %s" %row[0]
        print q
        cur.execute(q)
    q = "update row_desc set app_id = 1, mflag = 0;"
    print q
    cur.execute(q)
    conn.commit()

if __name__ == "__main__":
    truncate("mastodon", blade=True)
    reverseMarkAsDelete("diaspora", blade=False)
    truncateTableFromStencil("migration_registration")
    truncateTableFromStencil("evaluation")
    truncateTableFromStencil("display_flags")
    truncateTableFromStencil("user_table")
    truncateTableFromStencil("data_bags")
    resetRowDesc()