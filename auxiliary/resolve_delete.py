from QueryResolver import QueryResolver

if __name__ == "__main__":

    app_name = "hacker news"
    QR       = QueryResolver(app_name)

    with open("delete.queries") as fh: 
        queries = fh.read().split('\n')
        for q in queries:
            QR.resolveDelete(q)
            print QR.getResolvedQueries()
            QR.sendToDB(False)
            # break
        
        