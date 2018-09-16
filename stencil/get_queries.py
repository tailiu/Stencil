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

def translateAttribute(attr, table):
    sql = "SELECT column_name, mapping  \
            FROM app_mappings \
            INNER JOIN app_tables on app_mappings.table_id = app_tables.PK \
            WHERE column_name = '{0}' and table_name = '{1}'".format(attr, table)
    CUR.execute(sql)
    return CUR.fetchone()

def translateBasicSelectQuery(query):
    queryLower = query.lower()

    attributesStart = queryLower.find("select") + len("select")
    attributesEnd = queryLower.find("from")
    attributes = query[attributesStart : attributesEnd].split(",")

    tablesStart = queryLower.find("from") + len("from")
    tablesEnd = queryLower.find("where")
    tables = query[tablesStart : tablesEnd].split(",")

    condsStart = queryLower.find("where") + len("where")
    conds = query[condsStart:].split("and")

    condList = []
    for cond in conds:
        cond1 = cond.split('=')
        condList.append(cond1[0].strip())

    for i in range(len(attributes)):
        attributes[i] = attributes[i].strip()

    for i in range(len(condList)):
        condList[i] = condList[i].strip()

    attrList = list(set(attributes).union(condList))

    attributePairs = []
    for attr in attrList:
        translatedAttr = translateAttribute(attr, tables[0].strip()) # only for one table now
        attributePairs.append(translatedAttr)

    mergedReq = mergeRquests(attributePairs)

    return resolveRequest(mergedReq, query)

     
if __name__ == "__main__":

    CONN, CUR = getDBConn()

    sql = "SELECT By, Descendents, Id \
           FROM story  \
           WHERE By = \'\"Impossible\"\' and Id = 13075839" 
    
    translatedQuery = translateBasicSelectQuery(sql)
    
    print translatedQuery

    # sql1 = "SELECT PK, column_name \
    #         FROM app_mappings  \
    #         WHERE app_id = 3 and table_id = 1"

    CUR.execute(translatedQuery)

    for row in CUR.fetchall():
        print row