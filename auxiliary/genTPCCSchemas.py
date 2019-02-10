import random, copy, json
import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
MAX_COLS = 50

def getDBConn(db):
    conn = psycopg2.connect(dbname=db, user="root", host="10.224.45.158", port="26257")
    conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT) 
    return conn, conn.cursor()

def createOrGetDB(db, populate = False):
    conn, cur = getDBConn("tpcc")

    try:
        print "Creating database:", db
        cur.execute("CREATE DATABASE %s;"%db)
        cur.close(); conn.close()
        conn, cur = getDBConn(db)
        print "Importing TPCC Schema"
        cur.execute(open("./SQLS/tpcc_schema.sql", "r").read())
        addRandomTables(conn,cur, populate)
        print "DB: '%s', created!" % db        
    except Exception as e:
        print e
        # print "Error creating DB: '%s'! Already exists?" % db
    
    cur.close(); conn.close()

    return getDBConn(db)

def addRandomTables(conn, cur, populate = False):

    cur.execute("SHOW TABLES;")

    # tables  = [x[0] for x in cur.fetchall()]
    tables = ["warehouse", "district", "customer", "history", "orderr", "new_order", "item", "stock", "order_line"]
    columns = {}

    for table in tables:
        print "Adding random attributes to:", table
        cur.execute("SHOW COLUMNS FROM "+table)
        columns[table] = [x[0] for x in cur.fetchall()]
        numcols = random.randint(len(columns[table])+1, MAX_COLS)
        alter_table_sql = "ALTER TABLE %s "%table
        alter_col_sqls = []
        for i in range(len(columns[table]), numcols):
            # print "col%s"%i
            alter_table_sql += "ADD COLUMN col%s text, "%i
            alter_col_sql    = "ALTER TABLE %s ALTER COLUMN col%s SET DEFAULT md5(random()::text);"%(table, i)
            alter_col_sqls.append(alter_col_sql)
        alter_table_sql = alter_table_sql.strip(' ,') + ";"
        cur.execute(alter_table_sql)
        for asql in alter_col_sqls:
            cur.execute(asql)
        if populate:
            print "Populating TPCC data for table:", table
            d = "INSERT INTO"
            sqls = [d+e for e in open("./SQLS/%s.sql"%table, "r").read().split(d) if e]
            wrong_sqls = 0
            for sql in sqls:
                if sql.strip() != "INSERT INTO":
                    cur.execute(sql)
                else:
                    wrong_sqls += 1    
            # return
            # cur.execute(open("./SQLS/%s.sql"%table, "r").read())
            print table, "populated, wrong sqls:", wrong_sqls

def populateData(cur):
    tables = ["warehouse", "district", "customer", "history", "orderr", "new_order", "item", "stock", "order_line"]
    for table in tables:
        print "Truncate existing data:", table
        cur.execute("TRUNCATE TABLE "+table)
        print "Populating TPCC data for table:", table
        d = "INSERT INTO"
        sqls = [d+e for e in open("./SQLS/%s.sql"%table, "r").read().split(d) if e]
        wrong_sqls = 0
        for sql in sqls:
            if sql.strip() != "INSERT INTO":
                cur.execute(sql)
            else:
                wrong_sqls += 1    
        # return
        # cur.execute(open("./SQLS/%s.sql"%table, "r").read())
        print table, "populated, wrong sqls:", wrong_sqls

def printColDict(tables):
    for table in tables.keys():
        print "Table:- ", table
        print "Cols:- ",
        for col in tables[table]:
            print col,
        print

def genMapping(app1, app2):
    conn1, cur1 = createOrGetDB(app1, True)
    conn2, cur2 = createOrGetDB(app2)

    cur1.execute("SHOW TABLES;")

    mappings = {}

    for table in cur1.fetchall():
        table = table[0]
        mappings[table] = {table:{}}
        showtables = "SHOW COLUMNS FROM "+table
        cur1.execute(showtables); cur2.execute(showtables)
        cols1 = [x[0] for x in cur1.fetchall()]
        cols2 = [x[0] for x in cur2.fetchall()]
        mappings[table][table] = {col:col for col in cols1 if col in cols2}
    
    cur1.close(); conn1.close()
    cur2.close(); conn2.close()

    return mappings

def pushSchemaInfoIntoStencil(scur, app):
    conn, cur = getDBConn(app)

    sql = "INSERT INTO apps (app_name) VALUES ('%s') RETURNING row_id" % app
    scur.execute(sql)
    app_id = scur.fetchone()[0]
    # app_id = 1

    cur.execute("SHOW TABLES;")

    for table in cur.fetchall():
        table = table[0]
        showtables = "SHOW COLUMNS FROM "+table
        cur.execute(showtables)
        cols = ''.join(["(%s,'%s','%s','%s'),"%(app_id, table, x[0], x[1]) for x in cur.fetchall()])
        sql = "INSERT INTO app_schemas (app_id, table_name, column_name, data_type) VALUES "+cols.strip(", ")
        print sql
        scur.execute(sql)
        # print sql
        # print "\n\n\n"
        # break

def createMappingsInStencil(scur, apps):
    for app1 in apps:
        for app2 in apps:
            if app1 != app2:
                sql = """
                    INSERT INTO schema_mappings (source_attribute, dest_attribute)
                    SELECT as1.row_id, as2.row_id
                    FROM app_schemas as1 
                    JOIN app_schemas as2 ON as1.table_name = as2.table_name AND as1.column_name = as2.column_name
                    JOIN apps a1 ON a1.row_id = as1.app_id
                    JOIN apps a2 ON a2.row_id = as2.app_id
                    WHERE a1.app_name = '%s' AND a2.app_name='%s'
                """ % (app1, app2)
                scur.execute(sql)



if __name__ == "__main__" :

    app1 = "app1"
    apps = []
    total_apps = 10

    for i in range(1, total_apps):
        app = "app%s"%i
        apps.append(app)

    # app = "app9"
    # conn, cur = getDBConn(app)
    # print "Populating ", app
    # populateData(cur)

    # settings = {
    #     "user_table": "Customer",
    #     "key_column": "c_id",
    #     "mappings": {}
    # }
    # stencildb, stencilcur = getDBConn("stencil")
    # for app in apps:
    #     pushSchemaInfoIntoStencil(stencilcur, app)
    
    # createMappingsInStencil(stencilcur, apps)

    # with open('%s.json'%app1, 'w') as fp:
    #     for app in apps:
    #         settings["mappings"][app] = genMapping(app1, app)
    #     json.dump(settings, fp, indent=4)

    