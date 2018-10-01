import MySQLdb, json, os, re
from db import DB
from QueryResolver import QueryResolver

def getAppSchema(app_name):

    db = DB(host="10.224.45.162", user="zainmac", passwd="123")
    
    sql = """ SELECT    LOWER(app_tables.table_name), 
                        LOWER(app_schemas.column_name)
                FROM 	app_schemas 
                JOIN 	app_tables ON app_schemas.table_id = app_tables.PK
                JOIN 	apps ON app_tables.app_id = apps.PK
                WHERE   apps.app_name = "%s" """ % (app_name) 

    db.cursor.execute(sql)
    rows = db.cursor.fetchall()
    db.close()
    result = {}
    for row in rows:
        if row[0] in result.keys(): result[row[0]].append(row[1])
        else: result[row[0]] = [row[1]]
    return result

def getDatasetsPaths():

    twit_dir = "./datasets/Twitter/"
    fpaths = []

    for dir in os.listdir(twit_dir):
        dirpath = twit_dir + dir
        if os.path.isdir(dirpath):
            for file in os.listdir(dirpath):
                fpath = "%s/%s" % (dirpath, file)
                if fpath.find(".json") >= 0:
                    fpaths.append(fpath)
    return fpaths

def escapeString(item):
    # return str(item)
    # return 
    return MySQLdb.escape_string(json.dumps(item))

def getJSONDataFromFile(fpath):
    data = []
    with open(fpaths[0], "rb") as fh: 
        rows = fh.read().split("\n")
        for row in rows[:-1]:
            data.append(json.loads(row, encoding='utf-8'))
    return data


if __name__ == "__main__":

    app_name = "twitter"
    print "Started: %s" % app_name
    # db       = DB(host="10.224.45.162", user="zainmac", passwd="123", db=app_name)
    print "connected db"
    schema   = getAppSchema(app_name)
    print "got schema"
    fpaths   = getDatasetsPaths()
    print "got paths"
    QR       = QueryResolver(app_name)
    print "got QR"
    log_file = open("./twitter_log.txt", "wb")
    print "created Log"

    # db.truncateTables(["tweet", "user"])

    for fpath in fpaths[0:10]:
        print "Porting: ", fpath
        data = getJSONDataFromFile(fpath)
        
        for datum in data:
            
            t_attrs = [ x.lower() for x in datum.keys() if x.lower() in schema["tweet"] and x.lower() != "user"]
            if len(t_attrs) > 0:
                t_cols = "`%s`" % '`,`'.join(t_attrs)
                t_vals = "'%s'" % "','".join( [escapeString(datum[attr]) for attr in t_attrs] )

                if "user" in datum.keys():
                    u_attrs = [ x.lower() for x in datum["user"].keys() if x.lower() in schema["user"]]
                    if len(t_attrs) > 0:
                        u_cols = "`%s`" % '`,`'.join(u_attrs)
                        u_vals = "'%s'" % "','".join( [escapeString(datum["user"][attr]) for attr in u_attrs] )
                        u_sql  = "INSERT INTO user (%s) VALUES (%s);" % (u_cols, u_vals)
                        t_sql  = "INSERT INTO tweet (%s, `user`) VALUES (%s, '%s');" % (t_cols, t_vals, datum["user"]["id"])
                        # db.cursor.execute(u_sql)
                        # QR.resolveInsert(u_sql)
                        # QR.runQuery()
                else:
                    t_sql  = "INSERT INTO tweet (%s) VALUES (%s);" % (t_cols, t_vals)
                try:
                    # db.cursor.execute(t_sql)
                    QR.resolveInsert(t_sql)
                    QR.runQuery()
                except Exception as e:
                    print "------------------------------------------"
                    print "ERROR: %s" %e
                    print "------------------------------------------"
                    print t_sql
                    print QR.getResolvedQueries()
                    log_file.write(t_sql) 
                    log_file.write(str(QR.getResolvedQueries())) 
                    # break
        # db.conn.commit()
        QR.DBCommit()
