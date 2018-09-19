from QueryResolver import QueryResolver

if __name__ == "__main__":

    app_name = "hacker news"
    QR       = QueryResolver(app_name)

<<<<<<< HEAD
    with open("./datasets/hn_log.queries") as fh: 
=======
    with open("./dataset/hn_log.queries") as fh: 
>>>>>>> 162f34a9095d3da3cf843893d5f51a03bed5e9fb
        queries = fh.read().split('\n')
        for q in queries:
            QR.resolveInsert(q)
            print QR.getResolvedQueries()
            QR.sendToDB()
            print "------------------------------\n"
            # break
        
        
