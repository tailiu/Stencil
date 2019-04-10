import MySQLdb, json, os
from db import DB
from QueryResolver import QueryResolver

def getAppSchema(app_name):

    # db = DB(host="10.224.45.162", user="zainmac", passwd="123", port=3306)
    db = DB()
    
    sql = """ SELECT    LOWER(app_tables.table_name), 
                        LOWER(app_schemas.column_name)
                FROM 	app_schemas 
                JOIN 	app_tables ON app_schemas.table_id = app_tables.PK
                JOIN 	apps ON app_tables.app_id = apps.PK
                WHERE   apps.app_name = '%s' """ % (app_name) 

    db.cursor.execute(sql)
    rows = db.cursor.fetchall()
    db.close()
    result = {}
    for row in rows:
        if row[0] in result.keys(): result[row[0]].append(row[1])
        else: result[row[0]] = [row[1]]
    return result

def validString(item):

    item = json.dumps(item).strip('"[]')
    return MySQLdb.escape_string(item)


if __name__ == "__main__":

    app_name = "hacker news"
    db_name  = "hacker_news"
    # db       = DB(host="10.224.45.162", user="zainmac", passwd="123", db=db_name)
    db       = DB(db_name)
    hn_fpath = "./datasets/HackerNews/hn.full.json"
    # QR       = QueryResolver(app_name)
    log_file = open("./log.txt", "wb")

    # db.truncateTables(["story", "comment"])

    with open(hn_fpath) as fh: data = json.load(fh, encoding='utf-8')
    
    schema = getAppSchema(app_name)
    # hn_db  = DB(db_name)
    
    logical_queries = []
    i = 0
    for datum in data:
        i += 1
        if datum['type'].lower() in schema.keys(): # datum['type] here is by accident the table name..
            print i
            table_name = datum['type'] 
            
            attrs = [ x.lower() for x in datum.keys() if x.lower() in schema[table_name] ]
            
            cols = '"%s"' % '","'.join(attrs)
            vals = "E'%s'" % "',E'".join( [validString(datum[attr]) for attr in attrs] )
            
            sql  = "INSERT INTO %s (%s) VALUES (%s);" % (table_name, cols, vals)
            
            # print sql

            # logical_queries.append(sql)
            try:
                db.cursor.execute(sql)
            except Exception as e:
                print "**ERROR => ", e
                log_file.write("Logical DB Error: %s \n" % sql) 
                # break
            # try:
            #     QR.resolveInsert(sql)
            #     # print QR.getResolvedQueries()
            #     QR.runQuery()
            #     # QR.DBCommit()
            # except Exception as e:
            #     print "********ERROR in physical db: %s" % e
            #     log_file.write("Error: %s \n" % e) 
            #     log_file.write("Logical Query: %s \n" % sql) 
            #     log_file.write("Physical Queries: %s \n" % str(QR.getResolvedQueries())) 
            #     log_file.write("----------") 
            #     # break
    # QR.runAllQueries()
    # db.conn.commit()
    # QR.DBCommit()
    
    

    #################
    ## EXPORT TO FILE
    #################

    # hn_wpath = "./datasets/hn_log.queries"
    # with open(hn_wpath, "wb") as fh: 
    #     for q in logical_queries:
    #         fh.write("%s\n" % q)
    #     fh.seek(-1, os.SEEK_END)
    #     fh.truncate()
