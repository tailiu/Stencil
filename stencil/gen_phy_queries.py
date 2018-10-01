from QueryResolver import QueryResolver

if __name__ == "__main__":

    app_name = "hacker news"
    QR       = QueryResolver(app_name)

    with open("./datasets/hn_log.queries") as fh: 
        queries = fh.read().split('\n')
        start, end = 0, 10000
        for q in queries[start:end]:
            QR.resolveInsert(q)
            # print QR.getResolvedQueries()
            QR.runQuery()
    QR.DBCommit()
        
        
