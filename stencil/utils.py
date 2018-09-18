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
    if len(suppAttributeList) != 0: suppAttributes = findSuppTables(app_name, table, suppAttributeList)

    resolvedReq = resolveGetRowIDReq(table, baseAttributes, suppAttributes, conditions)

    # print resolvedReq
    
    CUR.execute(resolvedReq)
    return CUR.fetchall()

if __name__ == "__main__":
    CONN, CUR = getDBConn()
    
    result = getRowID(CUR, 'story', "By = 'Impossible' and Id = 13075839 or time = 1480550545")

    for row in result:
        print row