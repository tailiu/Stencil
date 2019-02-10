import MySQLdb

def getDBConn():
    db_conn = MySQLdb.connect(
        host   = "10.224.45.162",
        port   = 3306,
        user   = "freedom_tai",
        passwd = "123",
        db     = "stencil_storage",
    )
    return db_conn, db_conn.cursor()

if __name__ == "__main__":
    CONN, CUR = getDBConn()
    
    sql = "SELECT * FROM app_schemas"
    
    CUR.execute(sql)

    for row in CUR.fetchall():
        print row

    CUR.close()
    CONN.close() 