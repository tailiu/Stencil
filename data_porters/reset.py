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

def DeleteRowsFromMigrationRegistration(arg):
    conn, cur = getDB("stencil", blade=False)
    if arg == "log":
        is_log = "true"
    else:
        is_log = "false"
    q = "DELETE FROM migration_registration WHERE is_logical = "+is_log
    print q
    cur.execute(q)
    conn.commit()

def resetRowDesc():
    conn, cur = getDB("stencil", blade=False)
    
    q = "drop table migration_table;"
    print q
    cur.execute(q)
    conn.commit()
    q = "create table migration_table as table migration_table_backup;"
    print q
    conn.commit()
    cur.execute(q)
    constraints = [ "CREATE INDEX migration_table_app ON public.migration_table (app_id int4_ops,table_id int8_ops,mflag int4_ops,mark_as_delete bool_ops);",
                    "CREATE INDEX migration_table_app_table_rowid ON public.migration_table (app_id int4_ops,table_id int8_ops,row_id int8_ops);",
                    "CREATE INDEX migration_table_bag ON public.migration_table (bag bool_ops);",
                    "CREATE INDEX migration_table_dst_app ON public.migration_table (app_id int4_ops);",
                    "CREATE INDEX migration_table_dst_app_dst_rowid_org_rowid ON public.migration_table (app_id int4_ops,row_id int8_ops,group_id int8_ops);",
                    "CREATE INDEX migration_table_dst_rowid ON public.migration_table (row_id int8_ops);",
                    "CREATE INDEX migration_table_dst_rowid_org_rowid ON public.migration_table (row_id int8_ops,group_id int8_ops);",
                    "CREATE INDEX migration_table_dst_table ON public.migration_table (table_id int8_ops);",
                    "CREATE INDEX migration_table_dst_table_dst_app_dst_rowid_org_rowid ON public.migration_table (app_id int4_ops,row_id int8_ops,group_id int8_ops,table_id int8_ops);",
                    "CREATE INDEX migration_table_mark_as_deleted ON public.migration_table (mark_as_delete bool_ops);",
                    "CREATE INDEX migration_table_mflag ON public.migration_table (mflag int4_ops);",
                    "CREATE INDEX migration_table_migration_id ON public.migration_table (migration_id int4_ops);",
                    "CREATE INDEX migration_table_org_rowid ON public.migration_table (group_id int8_ops);",
                    "CREATE INDEX migration_table_user_id ON public.migration_table (user_id int4_ops);",
                    "ALTER TABLE public.migration_table ADD CONSTRAINT migration_table_pk PRIMARY KEY (app_id, table_id, group_id, row_id, mark_as_delete);",
                    "ALTER TABLE ONLY public.migration_table ALTER COLUMN mark_as_delete SET DEFAULT false;",
                    "ALTER TABLE ONLY public.migration_table ALTER COLUMN bag SET DEFAULT false;",
                    "ALTER TABLE ONLY public.migration_table ALTER COLUMN copy_on_write SET DEFAULT false;",
                    "ALTER TABLE ONLY public.migration_table ALTER COLUMN mflag SET DEFAULT 0;",
                    "ALTER TABLE ONLY public.migration_table ALTER COLUMN updated_at SET DEFAULT now();",
                    "ALTER TABLE ONLY public.migration_table ALTER COLUMN created_at SET DEFAULT now();",
                    "CREATE TRIGGER update_migration_table_changetimestamp BEFORE UPDATE ON migration_table FOR EACH ROW EXECUTE PROCEDURE update_changetimestamp_column();"]
    for q in constraints:
        print q
        cur.execute(q)
    conn.commit()

if __name__ == "__main__":
    if len(sys.argv) <= 1:
        print "provide an argument (phy, log, row, all, both), exiting."
    else:
        arg = sys.argv[1]
        if arg in ["phy", "all"]:
            truncatePhysicalTables()
        if arg in ["log", "row", "all", "both"]:
            # DeleteRowsFromMigrationRegistration(sys.argv[1])
            truncateTableFromStencil("migration_registration")
            truncateTableFromStencil("evaluation")
            truncateTableFromStencil("user_table")
            truncateTableFromStencil("display_flags")
            truncateTableFromStencil("txn_logs")
            if arg in ["log", "all", "both"]:
                truncate("mastodon", blade=True)
                reverseMarkAsDelete("diaspora", blade=False)
            if arg in ["row", "all", "both"]:
                # truncateTableFromStencil("data_bags")
                resetRowDesc()