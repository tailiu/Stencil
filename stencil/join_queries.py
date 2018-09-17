import MySQLdb

def getDBConn():
    db_conn = MySQLdb.connect(
        host   = "127.0.0.1",
        port   = 3307,
        user   = "root",
        passwd = "",
        db     = "stencil_storage",
    )
    return db_conn, db_conn.cursor()

def translateJoinQuery(originalQuery):
    return originalQuery

if __name__ == "__main__":

    CONN, CUR = getDBConn()

    sql = "select user, descendents, id\
            from base_1 join supplementary_1 on base_1.row_id  = supplementary_1.row_id \
            where user = 'impossible' and id = 13075839" 

    translatedQuery = translateJoinQuery(sql)

    CUR.execute(translatedQuery)

    for row in CUR.fetchall():
        print row