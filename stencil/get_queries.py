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

def translateAttribute(attr, table):
    sql = "SELECT column_name, mapping  \
            FROM app_mappings \
            INNER JOIN app_tables on app_mappings.table_id = app_tables.PK \
            WHERE column_name = '{0}' and table_name = '{1}'".format(attr, table)
    CUR.execute(sql)
    return CUR.fetchone()

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

def resolveRequest(req):
    print req
    # for k, v in req.iteritems():
    #     print k


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

    attributePairs = []
    for attr in attributes:
        translatedAttr = translateAttribute(attr.strip(), tables[0].strip())
        translatedAttr
        attributePairs.append(translatedAttr)
    
    print attributePairs
    mergedReq = mergeRquests(attributePairs)

    resolveRequest(mergedReq)

    return query

if __name__ == "__main__":

    CONN, CUR = getDBConn()

    sql = "SELECT By, Descendents, Id \
           FROM story  \
           WHERE By = 'Alice' and Id = 12" 
    
    translateBasicSelectQuery(sql)
    
    # sql1 = "SELECT PK, column_name \
    #         FROM app_mappings  \
    #         WHERE app_id = 3 and table_id = 1"

    # CUR.execute(sql1)

    # for row in CUR.fetchall():
    #     print row