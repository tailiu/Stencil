import MySQLdb, json

def getDBConn():
    db_conn = MySQLdb.connect(
        host   = "127.0.0.1",
        port   = 3307,
        user   = "root",
        passwd = "",
        db     = "stencil_storage",
    )
    return db_conn, db_conn.cursor()

def getAppSchema(app_name):
    
    sql = """ SELECT    LOWER(app_schemas.table_name), 
                        LOWER(app_schemas.column_name)
                FROM 	app_schemas 
                JOIN 	apps ON app_schemas.app_id = apps.PK
                WHERE   apps.app_name = "%s" """ % (app_name) 

    CUR.execute(sql)
    rows = CUR.fetchall()
    result = {}
    for row in rows:
        if row[0] in result.keys(): result[row[0]].append(row[1])
        else: result[row[0]] = [row[1]]
    return result

################ DB Globals ##
CONN, CUR = getDBConn()
##############################

if __name__ == "__main__":

    app_name = "hacker news"
    hn_fpath = "./data.json"

    with open(hn_fpath) as fh: data = json.load(fh, encoding='utf-8')
    
    schema = getAppSchema(app_name)
    
    logical_queries = []
    
    for datum in data:
        if datum['type'].lower() in schema.keys(): # datum['type] here is by accident the table name..
            table_name = datum['type'] 
            attrs = [ x.lower() for x in datum.keys() if x.lower() in schema[table_name] ]
            sql = "INSERT INTO %s ( " % table_name \
                  + ','.join(attrs) \
                  + " ) VALUES ( " \
                  + ','.join( ['"' + json.dumps(datum[attr]).strip('"[]') + '"' for attr in attrs] ) \
                  + " )"
            logical_queries.append(sql)

    hn_wpath = "hn_log.queries"
    with open(hn_wpath, "wb") as fh: 
        for q in logical_queries:
            fh.write("%s\n" % q)