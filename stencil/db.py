import MySQLdb
import psycopg2

class DB:
    # def __init__(self, db = "stencil_storage", host="127.0.0.1", user="root", passwd="", port=3306):
    def __init__(self, db = "stencil_storage"):
        # self.conn = MySQLdb.connect(
        #     host   = host,
        #     port   = port,
        #     user   = user,
        #     passwd = passwd,
        #     db     = db,
        # )
        self.conn = psycopg2.connect(
            database=db,
            user='root',
            sslmode='disable',
            port=26259,
            host='10.224.45.162'
        )
        self.conn.set_session(autocommit=True)
        self.cursor = self.conn.cursor()

    def truncateTables(self, tables):
        if len(tables) > 0:
            for table in tables:
                sql = "TRUNCATE %s" % table
                print sql
                self.cursor.execute(sql)
            self.conn.commit()
    
    def close(self):
        self.cursor.close()
        self.conn.close() 