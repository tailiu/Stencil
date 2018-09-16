import MySQLdb

def getDBConn():
    db_conn = MySQLdb.connect(
        host   = "127.0.0.1",
        port   = 3307,
        user   = "root",
        passwd = "",
        db     = "stencil_storage",
    )
    return db_conn, db_conn.cursor()

def mergeRquests(attributePairs):
    requestDict = {}
    for attrPair in attributePairs:
        tableName = attrPair[1].split('.')[0]
        attributeName = attrPair[1].split('.')[1]
        t = (attrPair[0], attributeName)
        if  tableName in requestDict:
            requestDict[tableName].append(t)
        else:
            requestDict[tableName] = []
            requestDict[tableName].append(t)
    return requestDict

def resolveRequest(req, originalQuery):
    resolvedQuery = originalQuery

    tables = []

    for k, v in req.iteritems():
        tables.append(k)
        for tuple in v:
            resolvedQuery = resolvedQuery.replace(tuple[0], tuple[1])

    fromWords = ''

    for i in range(len(tables)):
        if i != len(tables) - 1:
            fromWords += tables[i] + ' join '
        else:
            fromWords += tables[i]

    for i in range(len(tables)):
        if i == 0:
            fromWords += ' on ' + tables[i] + '.row_id = '
        else:
            fromWords += tables[i] + '.row_id'

    tablesStart = resolvedQuery.lower().find("from") + len("from")
    tablesEnd = resolvedQuery.lower().find("where")
    resolvedQuery = resolvedQuery[:tablesStart] +  ' ' + fromWords + ' ' + resolvedQuery[tablesEnd:]

    return resolvedQuery

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

def translateAttributesToBaseTables(app_name, table, attrList):
    attrStr = '('
    for i in range(len(attrList)):
        if i == len(attrList) - 1: attrStr += 'app_schemas.column_name = \'' + attrList[i] + '\')'
        else: attrStr += 'app_schemas.column_name = \'' + attrList[i] + '\' or '
        
    sql = "SELECT app_schemas.column_name, base_table_attributes.table_name, base_table_attributes.column_name\
            FROM base_table_attributes INNER JOIN physical_mappings INNER JOIN app_schemas INNER JOIN app_tables INNER JOIN apps\
            on base_table_attributes.PK = physical_mappings.physical_attribute \
            and app_schemas.PK = physical_mappings.logical_attribute\
            and app_tables.PK = app_schemas.table_id\
            and apps.PK = app_tables.app_id \
            WHERE app_name = '{0}' and app_tables.table_name = '{1}' and {2}".format(app_name, table, attrStr)

    CUR.execute(sql)
    return CUR.fetchall()
    
def translateBasicSelectQuery(query):
    query = query.lower()
    
    attributes = findBetweenStrings(query, 'select', 'from').split(',')
    tables = findBetweenStrings(query, 'from', 'where').split(',')
    condList = processConditions(findBetweenStrings(query, 'where', None))

    attributes = removeSpace(attributes)
    tables = removeSpace(tables)
    condList = removeSpace(condList)

    attrList = list(set(attributes).union(condList))

    print attrList

    baseAttributes = translateAttributesToBaseTables('hacker news', tables[0].strip(), attrList)
    suppAttributes = []

    print baseAttributes

    for attr in attrList:
        for var in baseAttributes:
            if attr != var[0].lower():
                continue
            suppAttributes.append(var[0])

    print suppAttributes

    # attributePairs = []
    # for attr in attrList:
    #     translatedAttr = translateAttributes('hacker news', tables[0].strip(), attr) # only for one table now
    #     attributePairs.append(translatedAttr)


    # print attributePairs
    # mergedReq = mergeRquests(attributePairs)

    # return resolveRequest(mergedReq, query)

     
if __name__ == "__main__":

    CONN, CUR = getDBConn()

    sql = "SELECT By, Descendents, Id \
           FROM story  \
           WHERE By = \'\"Impossible\"\' and Id = 13075839" 
    
    translatedQuery = translateBasicSelectQuery(sql)

    # sql1 = "SELECT PK, column_name \
    #         FROM app_mappings  \
    #         WHERE app_id = 3 and table_id = 1"

    # CUR.execute(translatedQuery)

    # for row in CUR.fetchall():
    #     print row