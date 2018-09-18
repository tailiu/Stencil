import MySQLdb
import re

def getDBConn():
    db_conn = MySQLdb.connect(
        host   = "127.0.0.1",
        port   = 3307,
        user   = "root",
        passwd = "",
        db     = "stencil_storage",
    )
    return db_conn, db_conn.cursor()

def findBetweenStrings(originalStr, str1, str2):
    strStart = originalStr.find(str1) + len(str1)
    strEnd = -1
    if str2 != None: strEnd = originalStr.find(str2)
    return originalStr[strStart : strEnd]

def processConditions(conds):
    conds = conds.split('and')

    condList = []
    for cond in conds:
        cond1 = cond.split('=')
        condList.append(cond1[0].strip())
    
    return condList

def removeSpace(l):
    for i in range(len(l)):
        l[i] = l[i].strip()
    return l
    
def formAttrStr(attrList):
    attrStr = '('
    for i in range(len(attrList)):
        if i == len(attrList) - 1: attrStr += 'app_schemas.column_name = \'' + attrList[i] + '\')'
        else: attrStr += 'app_schemas.column_name = \'' + attrList[i] + '\' or '
    return attrStr

def translateAttributesToBaseTables(app_name, table, attrList):
    attrStr = formAttrStr(attrList)
        
    sql = "SELECT app_schemas.column_name, base_table_attributes.table_name, base_table_attributes.column_name\
            FROM base_table_attributes INNER JOIN physical_mappings INNER JOIN app_schemas INNER JOIN app_tables INNER JOIN apps\
            on base_table_attributes.PK = physical_mappings.physical_attribute\
            and app_schemas.PK = physical_mappings.logical_attribute\
            and app_tables.PK = app_schemas.table_id\
            and apps.PK = app_tables.app_id \
            WHERE app_name = '{0}' and app_tables.table_name = '{1}' and {2}".format(app_name, table, attrStr)

    CUR.execute(sql)
    return CUR.fetchall()

def findSuppTables(app_name, table, suppAttributes):
    attrStr = formAttrStr(suppAttributes)

    sql = "SELECT app_schemas.column_name, supplementary_tables.supplementary_table\
            FROM supplementary_tables INNER JOIN app_schemas INNER JOIN app_tables INNER JOIN apps\
            on supplementary_tables.table_id = app_schemas.table_id\
            and app_tables.PK = app_schemas.table_id\
            and apps.PK = app_tables.app_id \
            WHERE app_name = '{0}' and app_tables.table_name = '{1}' and {2}".format(app_name, table, attrStr)
    
    CUR.execute(sql)
    return CUR.fetchall()

def resolveRequest(query, baseAttributes, suppAttributes):
    for attr in baseAttributes: query = re.sub(r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[2].lower(), query)
    for attr in suppAttributes: query = re.sub(r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[0].lower(), query)

    tables = set()
    for attr in baseAttributes: tables.add(attr[1])
    for attr in suppAttributes: tables.add(attr[1])

    fromTables = ''
    for table in tables: fromTables += ' join ' + table
    fromTables = fromTables.replace('join', '', 1)
    fromTables += ' on '
    for table in tables: fromTables += ' = ' + table + '.row_id '
    fromTables = fromTables.replace('=', '', 1)

    tablesStart = query.find("from") + len("from")
    tablesEnd = query.find("where")
    query = query[:tablesStart] +  ' ' + fromTables + ' ' + query[tablesEnd:]

    return query

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

def translateBasicSelectQuery(query):
    query = query.lower()
    
    table = findBetweenStrings(query, 'from', 'where').strip() # Assume there is only one table

    if query.find('*') == -1:
        attributes = findBetweenStrings(query, 'select', 'from').split(',')
        attributes = removeSpace(attributes)
    else: 
        attributes = findAllAttributes('hacker news', table)
        attrStr = ''
        for attr in attributes: attrStr += ", " + attr
        attrStr = attrStr.replace(',', '', 1)
        query = query.replace('*', attrStr)
        query = query.lower()

    condList = processConditions(findBetweenStrings(query, 'where', None))
    condList = removeSpace(condList)

    attrList = list(set(attributes).union(condList))

    baseAttributes = translateAttributesToBaseTables('hacker news', table, attrList)

    suppAttributeList = []
    for attr in attrList:
        find = False
        for var in baseAttributes:
            if attr.lower() == var[0].lower(): 
                find = True
                break
        if not find: suppAttributeList.append(attr)

    suppAttributes = findSuppTables('hacker news', table, suppAttributeList)

    return resolveRequest(query, baseAttributes, suppAttributes)
     
if __name__ == "__main__":

    CONN, CUR = getDBConn()

    sql = "SELECT By, Descendents, Id, Retrieved_on, Score\
           FROM story  \
           WHERE By = 'Impossible' and Id = 13075839" 
    
    sql1 = "SELECT * \
            FROM story  \
            WHERE By = 'Impossible' and Id = 13075839"

    sql2 = "SELECT * \
            FROM comment  \
            WHERE Id = 13075839 + 1 "

    translatedQuery = translateBasicSelectQuery(sql2)
    print translatedQuery

    CUR.execute(translatedQuery)
    for row in CUR.fetchall():
        print row