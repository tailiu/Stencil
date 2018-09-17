import MySQLdb, json, sqlparse, uuid
from QueryResolver import QueryResolver

if __name__ == "__main__":

    app_name = "hacker news"
    QR       = QueryResolver(app_name)

    with open("hn_log.queries") as fh: 
        queries = fh.read().split('\n')
        for q in queries:
            QR.resolveInsert(q)
            QR.sendToDB()
            # break
        
        
