import psycopg2
import json

def getDBConn(app_name):
    conn     = psycopg2.connect(dbname=app_name, user="cow", password="123456" ,host="10.230.12.86", port="5432")
    cursor   = conn.cursor()
    return conn, cursor

stencilConn, stencilCursor = getDBConn("stencil")

def deletePhysicalTables():
    
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    stencilCursor.execute(tableq)
    for row in stencilCursor.fetchall():
        table = row[0]
        if table != "supplementary_tables" and ("supplementary_" in table or "base_" in table):
            q = 'DROP TABLE "%s" CASCADE;'%table
            print q
            stencilCursor.execute(q)
            stencilConn.commit()

def truncatePhysicalTables():
    tables = ["supplementary_tables", "physical_schema", "physical_mappings"]
    for table in tables:
        sql = 'TRUNCATE "%s" RESTART IDENTITY CASCADE;'%table
        print sql
        stencilCursor.execute(sql)
        stencilConn.commit()

def truncateAllTables():
    
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    stencilCursor.execute(tableq)
    for row in stencilCursor.fetchall():
        sql = 'TRUNCATE "%s" RESTART IDENTITY CASCADE;'%row[0]
        print sql
        stencilCursor.execute(sql)
        stencilConn.commit()

def truncateTable(table):

    sql = 'TRUNCATE "%s" RESTART IDENTITY CASCADE;'%table
    stencilCursor.execute(sql)
    stencilConn.commit()

def getTablesForApp(app_name):
    conn, cursor = getDBConn(app_name)
    tables_sql  = "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';"
    cursor.execute(tables_sql)
    table_list = []
    for trow in cursor.fetchall():
        table_list.append(trow[0])
    return table_list

def populateApps(apps):
    for i in range(0, len(apps)):
        pk, app = i+1, apps[i]
        sql = "INSERT INTO apps (app_name, pk) VALUES ('%s', '%d')" % (app, pk)
        print sql
        stencilCursor.execute(sql)
    stencilConn.commit()
    
def getAppID(app_name):
    
    sql  = "SELECT pk FROM apps WHERE app_name = '%s';"%app_name
    stencilCursor.execute(sql)
    row = stencilCursor.fetchone()
    return row[0]

def insertAppTables(app_name):
    app_id = getAppID(app_name) 
    sql = "INSERT INTO app_tables (app_id, table_name) VALUES "
    for table in getTablesForApp(app_name):
        sql += "('%s', '%s'), " % (app_id, table)
    sql = sql.strip(", ")
    print sql
    stencilCursor.execute(sql)
    stencilConn.commit()

def populateAppSchema(app_name):

    appConn, appCursor = getDBConn(app_name)

    table_sql = "SELECT app_tables.pk, app_tables.table_name FROM app_tables JOIN apps ON app_tables.app_id = apps.pk WHERE apps.app_name = '%s'"%app_name
    stencilCursor.execute(table_sql)

    for trow in stencilCursor.fetchall():
        table_id = trow[0]
        table_name = trow[1]
        # columns_sql = 'SHOW COLUMNS FROM "%s"' % table_name
        columns_sql = "select column_name, data_type from INFORMATION_SCHEMA.COLUMNS where table_name = '%s'" % table_name
        appCursor.execute(columns_sql)
        for crow in appCursor.fetchall():
            column_name = crow[0]
            data_type = crow[1]
            if data_type == "ARRAY":
                data_type = "text"
            isql = "INSERT INTO app_schemas (table_id, column_name, data_type) VALUES (%d, '%s', '%s')" % (table_id, column_name, data_type)
            print isql
            stencilCursor.execute(isql)
    stencilConn.commit()

def addSchemaMappings(app_name):

    app_id = getAppID(app_name) 
    with open('../stencil/config/app_settings/mappings.json', 'r') as mappingFile: file_data = mappingFile.read()
    mappings = json.loads(file_data)
    for appMapping in mappings["allMappings"]:
        if appMapping["fromApp"] == app_name:
            for mappingToApp in appMapping["toApps"]:
                mappedApp = mappingToApp["name"]
                mappings = mappingToApp["mappings"]
                for mapping in mappings:
                    mappedTables = mapping["toTables"]
                    for mappedTable in mappedTables:
                        mappedTableName = mappedTable["table"]
                        for mappedCol, mappedFromAttr in mappedTable["mapping"].items():
                            if "$" in mappedFromAttr or "#" in mappedFromAttr:
                                pass
                            else:
                                mappedFromAttr = mappedFromAttr.split(".")
                                mapperTable = mappedFromAttr[0]
                                mapperCol = mappedFromAttr[1]

                                sql = "select app_schemas.pk from app_schemas join app_tables on app_schemas.table_id = app_tables.pk join apps on apps.pk = app_tables.app_id where apps.app_name = '%s' and app_tables.table_name = '%s' and app_schemas.column_name = '%s'"
                                
                                _sql = sql%(app_name,mapperTable,mapperCol)
                                stencilCursor.execute(_sql)
                                mapperAttrID = stencilCursor.fetchone()[0]
                                
                                _sql = sql%(mappedApp,mappedTableName,mappedCol)
                                stencilCursor.execute(_sql)
                                try:
                                    mappedAttrID = stencilCursor.fetchone()[0]
                                except Exception as e:
                                    print "ERROR ENCOUNTERED! EXIT!"
                                    print e
                                    print _sql
                                    exit(0)

                                isql = "INSERT INTO schema_mappings (source_attribute, dest_attribute) VALUES (%d, %d)"
                                print app_name,mapperTable,mapperCol, mapperAttrID, "=>", mappedApp, mappedTableName, mappedCol, "id", mappedAttrID
                                stencilCursor.execute(isql%(mapperAttrID, mappedAttrID))                                
    stencilConn.commit()

if __name__ == "__main__":
    apps = ["diaspora", "mastodon","twitter"]
    deletePhysicalTables()
    truncateAllTables()
    populateApps(apps)
    for app_name in apps:
        insertAppTables(app_name)
        populateAppSchema(app_name)
    for app_name in apps:
        addSchemaMappings(app_name)