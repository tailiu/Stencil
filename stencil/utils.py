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

def findSuppTables(CUR, app_name, table, suppAttributes):
    attrStr = utils.formAttrStr(suppAttributes)

    sql = "SELECT app_schemas.column_name, supplementary_tables.supplementary_table\
            FROM supplementary_tables INNER JOIN app_schemas INNER JOIN app_tables INNER JOIN apps\
            on supplementary_tables.table_id = app_schemas.table_id\
            and app_tables.PK = app_schemas.table_id\
            and apps.PK = app_tables.app_id \
            WHERE app_name = '{0}' and app_tables.table_name = '{1}' and {2}".format(app_name, table, attrStr)
    
    CUR.execute(sql)
    return CUR.fetchall()

def findAllAttributes(CUR, app_name, table):
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

def findBetweenStrings(originalStr, str1, str2):
    strStart = originalStr.find(str1) + len(str1)
    strEnd = -1
    if str2 != None: strEnd = originalStr.find(str2)
    return originalStr[strStart : strEnd]

def formAttrStr(attrList):
    attrStr = '('
    for i in range(len(attrList)):
        if i == len(attrList) - 1: attrStr += 'app_schemas.column_name = \'' + attrList[i] + '\')'
        else: attrStr += 'app_schemas.column_name = \'' + attrList[i] + '\' or '
    return attrStr

def translateAttributesToBaseTables(CUR, app_name, table, attrList):
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

def processConditions(conds):
    conds = re.split('and | or', conds)

    condList = []
    for cond in conds:
        cond1 = cond.split('=')
        condList.append(cond1[0].strip())
    
    return condList

def removeSpace(l):
    for i in range(len(l)):
        l[i] = l[i].strip()
    return l

def resolveGetRowIDReq(table, baseAttributes, suppAttributes, conditions):
    for attr in baseAttributes: conditions = re.sub(r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[2].lower(), conditions)
    for attr in suppAttributes: conditions = re.sub(r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[0].lower(), conditions)

    tables = set()
    for attr in baseAttributes: tables.add(attr[1])
    for attr in suppAttributes: tables.add(attr[1])

    if len(tables) == 1:
        fromTables = tables.pop()
        oneTable = fromTables
    else:
        fromTables = ''
        for table in tables: 
            fromTables += ' join ' + table
            oneTable = table
        fromTables = fromTables.replace('join', '', 1)
        fromTables += ' on '
        for table in tables: fromTables += ' = ' + table + '.row_id '
        fromTables = fromTables.replace('=', '', 1)

    query = 'SELECT ' + oneTable + '.row_id ' + 'FROM ' + fromTables + ' WHERE ' + conditions
    return query

def getRowID(CUR, app_name, logicalTableName, conditions):

    conditions = conditions.lower()

    table = logicalTableName.strip()
    condList = processConditions(conditions)

    baseAttributes = translateAttributesToBaseTables(CUR, app_name, table, condList)

    suppAttributeList = []
    for attr in condList:
        find = False
        for var in baseAttributes:
            if attr.lower() == var[0].lower(): 
                find = True
                break
        if not find: suppAttributeList.append(attr)

    suppAttributes = ()
<<<<<<< HEAD
    if len(suppAttributeList) != 0: suppAttributes = findSuppTables(app_name, table, suppAttributeList)
=======
    if len(suppAttributeList) != 0: suppAttributes = findSuppTables(CUR, 'hacker news', table, suppAttributeList)
>>>>>>> 5ed43a91dcb52dd470a0543ebdc12bf3c6de56e8

    resolvedReq = resolveGetRowIDReq(table, baseAttributes, suppAttributes, conditions)

    # print resolvedReq
    
    CUR.execute(resolvedReq)
    return CUR.fetchall()

if __name__ == "__main__":
    CONN, CUR = getDBConn()
    
    result = getRowID(CUR, 'story', "By = 'Impossible' and Id = 13075839 or time = 1480550545")

    for row in result:
        print row