import psycopg2
import json

def populateAppSchema(app_name):

    appConn     = psycopg2.connect(dbname=app_name, user="root", host="10.230.12.75", port="26257")
    stencilConn = psycopg2.connect(dbname="stencil", user="root", host="10.230.12.75", port="26257")
    appCursor, stencilCursor = appConn.cursor(), stencilConn.cursor()

    table_sql = "SELECT app_tables.rowid, app_tables.table_name FROM app_tables JOIN apps ON app_tables.app_id = apps.rowid WHERE apps.app_name = '%s'"%app_name
    stencilCursor.execute(table_sql)

    for trow in stencilCursor.fetchall():
        table_id = trow[0]
        table_name = trow[1]
        columns_sql = 'SHOW COLUMNS FROM "%s"' % table_name
        appCursor.execute(columns_sql)
        for crow in appCursor.fetchall():
            column_name = crow[0]
            data_type = crow[1]
            isql = "INSERT INTO app_schemas (table_id, column_name, data_type) VALUES (%d, '%s', '%s')" % (table_id, column_name, data_type)
            stencilCursor.execute(isql)
    stencilConn.commit()

def addSchemaMappings(app_name):
    with open('./transaction/config/app_settings/mappings.json', 'r') as mappingFile: file_data = mappingFile.read()
    stencilConn = psycopg2.connect(dbname="stencil", user="root", host="10.230.12.75", port="26257")
    stencilCursor = stencilConn.cursor()
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
                            try:
                                mappedFromAttr = mappedFromAttr.split(".")
                                mapperTable = mappedFromAttr[0]
                                mapperCol = mappedFromAttr[1]

                                sql = "select app_schemas.rowid from app_schemas join app_tables on app_schemas.table_id = app_tables.rowid join apps on apps.rowid = app_tables.app_id where apps.app_name = '%s' and app_tables.table_name = '%s' and app_schemas.column_name = '%s'"
                                
                                _sql = sql%(app_name,mapperTable,mapperCol)
                                stencilCursor.execute(_sql)
                                mapperAttrID = stencilCursor.fetchone()[0]
                                
                                _sql = sql%(mappedApp,mappedTableName,mappedCol)
                                stencilCursor.execute(_sql)
                                mappedAttrID = stencilCursor.fetchone()[0]

                                isql = "INSERT INTO schema_mappings (app1_attribute, app2_attribute) VALUES (%d, %d)"
                                stencilCursor.execute(isql%(mapperAttrID, mappedAttrID))

                                print app_name,mapperTable,mapperCol, "id", mapperAttrID, mappedApp, mappedTableName, mappedCol, "id", mappedAttrID
                                print isql
                            except IndexError as e:
                                print mappedFromAttr, e
    stencilConn.commit()


if __name__ == "__main__":
    for app_name in ["diaspora", "mastodon"]:
        # populateAppSchema(app_name)
        addSchemaMappings(app_name)
        # break