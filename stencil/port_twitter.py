import MySQLdb, json, os
from db import DB

def getAppSchema(app_name):

    db = DB()
    
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

def validString(item):

    item = json.dumps(item).strip('"[]')
    return MySQLdb.escape_string(item)


if __name__ == "__main__":

    app_name = "hacker news"
    db_name  = "hacker_news"
    hn_fpath = "./data.json"
    hn_fpath = "/Users/zain/Documents/DataSets/HackerNews/hn.full.json"

    with open(hn_fpath) as fh: data = json.load(fh, encoding='utf-8')
    
    schema = getAppSchema(app_name)
    hn_db  = DB(db_name)
    
    logical_queries = []
    
    for datum in data:
        if datum['type'].lower() in schema.keys(): # datum['type] here is by accident the table name..
            
            table_name = datum['type'] 
            
            attrs = [ x.lower() for x in datum.keys() if x.lower() in schema[table_name] ]
            
            cols = "`%s`" % '`,`'.join(attrs)
            vals = "'%s'" % "','".join( [validString(datum[attr]) for attr in attrs] )
            
            sql  = "INSERT INTO %s (%s) VALUES (%s);" % (table_name, cols, vals)
            
            print sql

            logical_queries.append(sql)
            hn_db.cursor.execute(sql)
    hn_db.conn.commit()

    # hn_wpath = "hn_log.queries"
    # with open(hn_wpath, "wb") as fh: 
    #     for q in logical_queries:
    #         fh.write("%s\n" % q)
    #     fh.seek(-1, os.SEEK_END)
    #     fh.truncate()




