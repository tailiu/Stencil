from QueryResolver import QueryResolver

if __name__ == "__main__":

    app_name = "hacker news"
    QR       = QueryResolver(app_name)

    with open("./datasets/hn_log.queries") as fh: 
        queries = fh.read().split('\n')
        for q in queries:
            QR.resolveInsert(q)
            print QR.getResolvedQueries()
            QR.sendToDB()
            print "------------------------------\n"
            # break
        
        
