import MySQLdb
import utils

def getDBConn():
    db_conn = MySQLdb.connect(
        host   = "127.0.0.1",
        port   = 3307,
        user   = "root",
        passwd = "",
        db     = "stencil_storage",
    )
    return db_conn, db_conn.cursor()

def findAllAttributes(app_name, table):
    sql = "SELECT app_schemas.column_name\
        FROM app_schemas INNER JOIN app_tables INNER JOIN apps\
        on app_tables.PK = app_schemas.table_id\
        and apps.PK = app_tables.app_id \
        WHERE app_name = '{0}' and app_tables.table_name = '{1}'".format(app_name, table)

    CUR.execute(sql)
    attributes = CUR.fetchall()

    attrList = []
    for attr in attributes:
        attrList.append(attr[0])
    return attrList

def formAttrStr(attrList):
    attrStr = '('
    for i in range(len(attrList)):
        if i == len(attrList) - 1: attrStr += 'app_schemas.column_name = \'' + attrList[i] + '\')'
        else: attrStr += 'app_schemas.column_name = \'' + attrList[i] + '\' or '
    return attrStr

def formTableStr(tables):
    tableStr = '('
    for table in tables: tableStr += ' or app_tables.table_name = \'' + table + '\''
    tableStr += ')'
    return tableStr.replace('or', '', 1)

def translateAttributesToBaseTables(CUR, app_name, tables, attrList):
    attrStr = formAttrStr(attrList)
    tableStr = formTableStr(tables)

    # print tableStr
    sql = "SELECT app_schemas.column_name, base_table_attributes.table_name, base_table_attributes.column_name\
            FROM base_table_attributes INNER JOIN physical_mappings INNER JOIN app_schemas INNER JOIN app_tables INNER JOIN apps\
            on base_table_attributes.PK = physical_mappings.physical_attribute\
            and app_schemas.PK = physical_mappings.logical_attribute\
            and app_tables.PK = app_schemas.table_id\
            and apps.PK = app_tables.app_id \
            WHERE app_name = '{0}' and {1} and {2}".format(app_name, tableStr, attrStr)

    # print sql

    CUR.execute(sql)
    return CUR.fetchall()

def translateJoinQuery(CUR, query):
    query = query.lower()

    tableInfo = utils.findBetweenStrings(query, 'from', 'where').strip()
    tables = tableInfo[:tableInfo.find('on')].split('inner join')
    tables = utils.removeSpace(tables)
    
    if query.find('*') == -1:
        attributes = utils.findBetweenStrings(query, 'select', 'from').split(',')
        attributes = utils.removeSpace(attributes)
    else: 
        attributes = utils.findAllAttributes('hacker news', table)
        attrStr = ''
        for attr in attributes: attrStr += ", " + attr
        attrStr = attrStr.replace(',', '', 1)
        query = query.replace('*', attrStr)
        query = query.lower()

    condList = utils.processConditions(utils.findBetweenStrings(query, 'where', None))
    condList = utils.removeSpace(condList)

    attrList = list(set(attributes).union(condList))

    print tables

    baseAttributes = translateAttributesToBaseTables(CUR, 'hacker news', tables, attrList)

    # print baseAttributes
    # print attrList
    # print attributes
    # print tables
    # print condList


if __name__ == "__main__":

    CONN, CUR = getDBConn()

    sql = "SELECT descendents, kids, parent, story.id\
            FROM story INNER JOIN comment on story.id = comment.parent \
            WHERE story.By = 'edblarney'"

    translatedQuery = translateJoinQuery(CUR, sql)

    # CUR.execute(translatedQuery)

    # for row in CUR.fetchall():
    #     print row