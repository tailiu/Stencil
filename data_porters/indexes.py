import psycopg2

def getDBConn(db):
    conn = psycopg2.connect(dbname=db, user="cow", password="123456", host="10.230.12.86", port="5432")
    cursor = conn.cursor()
    return conn, cursor

conn, cur = getDBConn("stencil")

def getPhysicalTables():
    tables = []
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    cur.execute(tableq)
    for row in cur.fetchall():
        table = row[0]
        if table != "supplementary_tables" and ("supplementary_" in table or "base_" in table):
            tables.append(table)
    return tables

if __name__ == "__main__":
    for table in getPhysicalTables():
        colsql = "select column_name from INFORMATION_SCHEMA.COLUMNS where table_name = '%s'"%table
        cur.execute(colsql)
        idxsql = ""
        for row in cur.fetchall():
            column_name = row[0]
            if column_name != "app_id" and (column_name == "id" or "_id" in column_name):
                idxsql += "CREATE INDEX %s_%s_idx ON public.%s (%s); "%(table, column_name, table, column_name)
        if idxsql != "":
            print(idxsql)
            cur.execute(idxsql)
            conn.commit()