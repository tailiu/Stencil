import psycopg2

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

if __name__ == "__main__":
    for app_name in ["mastodon"]:
        populateAppSchema(app_name)