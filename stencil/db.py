import MySQLdb

class DB:

    def __init__(self):
        self.conn = MySQLdb.connect(
            host   = "127.0.0.1",
            port   = 3307,
            user   = "root",
            passwd = "",
            db     = "stencil_storage",
        )
        self.cursor = self.conn.cursor()
    
    def close(self):
        self.cursor.close()
        self.conn.close() 