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

def getSchemaMapping(app_name, table_name):
    sql =  "SELECT GROUP_CONCAT(LOWER(app_mappings.column_name), ',', LOWER(app_mappings.mapping)) \
            FROM app_mappings JOIN apps \
            ON app_mappings.app_id = apps.PK \
            JOIN  app_tables \
            ON app_mappings.table_id = app_tables.PK \
            WHERE apps.app_name = '%s' AND app_tables.table_name = '%s'" % (app_name, table_name)
    CUR.execute(sql)
    result = CUR.fetchone()[0].split(",")
    return {result[i]: result[i+1] for i in range(0, len(result), 2)}

def getTableNames(app_name):
    sql = "SELECT table_name \
           FROM app_tables \
           JOIN apps ON apps.PK = app_tables.app_id \
           WHERE apps.app_name = '%s'" % app_name
    CUR.execute(sql)
    return [x[0].lower() for x in CUR.fetchall()]

################ DB Globals ##
CONN, CUR = getDBConn()
##############################

if __name__ == "__main__":

    hn_fpath = "/Users/zain/Documents/DataSets/HackerNews/hn.min.json"

    with open(hn_fpath) as fh: data = json.load(fh, encoding='utf-8')
    
    app_name = "hacker news"
    schemas  = {}

    for table_name in getTableNames(app_name):
        schemas[table_name] = getSchemaMapping(app_name, table_name)

    logical_queries = []
    physical_queries = []
    
    for datum in data:
        if datum['type'].lower() in schemas.keys():
            table_name = datum['type'] 
            attrs = [ x.lower() for x in datum.keys() if x.lower() in schemas[table_name].keys() ]
            sql1 = "INSERT INTO %s ( " % table_name \
                  + ','.join(attrs) \
                  + " ) VALUES ( " \
                  + ','.join([json.dumps(datum[attr]) for attr in attrs]) \
                  + " )"
            sql2 = "INSERT INTO %s ( " % table_name \
                  + ','.join([schemas[table_name][attr] for attr in attrs]) \
                  + " ) VALUES ( " \
                  + ','.join([json.dumps(datum[attr]) for attr in attrs]) \
                  + " )"
            logical_queries.append(sql1)
            physical_queries.append(sql2)

    hn_wpath = "hn_log.queries"
    with open(hn_wpath, "wb") as fh: 
        for q in logical_queries:
            # print q
            fh.write("%s\n" % q)