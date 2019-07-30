import psycopg2

def getDBConn(db):
    conn = psycopg2.connect(dbname=db, user="cow", password="123456", host="10.230.12.86", port="5432")
    cursor = conn.cursor()
    return conn, cursor

conn, cur = getDBConn("stencil")

for i in range(771,860):
    colsql = "select column_name from INFORMATION_SCHEMA.COLUMNS where table_name = 'supplementary_%d'"%i
    cur.execute(colsql)
    idxsql = ""
    for row in cur.fetchall():
        column_name = row[0]
        if column_name == "id" or "_id" in column_name:
            idxsql += "CREATE INDEX supplementary_%d_%s_idx ON public.supplementary_%d (%s); "%(i, column_name,i, column_name)
    if idxsql != "":
        print(idxsql)
        cur.execute(idxsql)
        conn.commit()