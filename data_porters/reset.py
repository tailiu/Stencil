import psycopg2, datetime, time, sys

def getDB(dbname, blade=False):
    if blade:
        host = "10.230.12.75"
    else:
        host = "10.230.12.86"
    conn = psycopg2.connect(dbname=dbname, user="cow", host=host, port="5432", password="123456")
    cursor = conn.cursor()
    return conn, cursor

def truncatePhysicalTables():
    conn, cur = getDB("stencil", blade=False)
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    cur.execute(tableq)
    for row in cur.fetchall():
        table = row[0]
        if table != "supplementary_tables" and ("supplementary_" in table or "base_" in table):
            q = 'TRUNCATE "%s" RESTART IDENTITY CASCADE;'%table
            print q
            cur.execute(q)
            conn.commit()
    truncateTableFromStencil("migration_table")
    
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
    
    # q = "select rowid from row_desc group by rowid having count(*) > 1;"
    # cur.execute(q)
    # for row in cur.fetchall():
    #     q = "delete from row_desc where app_id != 1 and rowid = %s" %row[0]
    #     print q
    #     cur.execute(q)
    # q = "update row_desc set app_id = 1, mflag = 0;"
    q = "delete from migration_table where app_id != 1"
    print q
    cur.execute(q)
    q = "update migration_table set mark_as_delete = false, mflag = 0, bag = false, migration_id = NULL, user_id = NULL, copy_on_write = false;"
    print q
    cur.execute(q)
    conn.commit()

if __name__ == "__main__":
    if len(sys.argv) <= 1:
        print "provide an argument (phy, log, row, all), exiting."
    else:
        arg = sys.argv[1]
        if arg in ["phy", "all"]:
            truncatePhysicalTables()
        if arg in ["log", "row", "all"]:
            truncateTableFromStencil("migration_registration")
            truncateTableFromStencil("evaluation")
            truncateTableFromStencil("user_table")
            truncateTableFromStencil("display_flags")
            if arg in ["log", "all"]:
                truncate("mastodon", blade=True)
                reverseMarkAsDelete("diaspora", blade=False)
            if arg in ["row", "all"]:
                # truncateTableFromStencil("data_bags")
                resetRowDesc()