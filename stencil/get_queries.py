import MySQLdb
import re
import utils
import datetime

def resolveRequest(query, baseAttributes, suppAttributes):
    for attr in baseAttributes: query = re.sub(r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[2].lower(), query)
    for attr in suppAttributes: query = re.sub(r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[0].lower(), query)

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

    tablesStart = query.find("from") + len("from")
    tablesEnd = query.find("where")
    query = query[:tablesStart] +  ' ' + fromTables + ' ' + query[tablesEnd:]

    return query

def translateBasicSelectQuery(CUR, query):
    query = query.lower()
    
    table = utils.findBetweenStrings(query, 'from', 'where').strip() # Assume there is only one table

    if query.find('*') == -1:
        attributes = utils.findBetweenStrings(query, 'select', 'from').split(',')
        attributes = utils.removeSpace(attributes)
    else: 
        attributes = utils.findAllAttributes(CUR, 'twitter', table)
        attrStr = ''
        for attr in attributes: attrStr += ", " + attr
        attrStr = attrStr.replace(',', '', 1)
        query = query.replace('*', attrStr)
        query = query.lower()

    condList = utils.processConditions(utils.findBetweenStrings(query, 'where', None))
    condList = utils.removeSpace(condList)
    
    attrList = list(set(attributes).union(condList))

    baseAttributes = utils.translateAttributesToBaseTables(CUR, 'twitter', table, attrList)
    
    suppAttributeList = []
    for attr in attrList:
        find = False
        for var in baseAttributes:
            if attr.lower() == var[0].lower(): 
                find = True
                break
        if not find: suppAttributeList.append(attr)

    suppAttributes = ()
    if len(suppAttributeList) != 0: suppAttributes = utils.findSuppTables(CUR, 'twitter', table, suppAttributeList)

    return resolveRequest(query, baseAttributes, suppAttributes)


if __name__ == "__main__":

    CONN, CUR = utils.getDBConn()

    sql = "SELECT By, Descendents, Id, Retrieved_on, Score\
           FROM story  \
           WHERE By = 'Impossible' and Id = 13075839" 
    
    sql1 = "SELECT * \
            FROM story  \
            WHERE By = 'Impossible' and Id = 13075839"

    sql2 = "SELECT * \
            FROM comment  \
            WHERE Id = 13075839 + 1 "


    sql2 = "SELECT * FROM tweet where user =2238942602"

    translatedQuery = translateBasicSelectQuery(CUR, sql2)
    print "translatedQuery: ", translatedQuery
    pre_time = datetime.datetime.now().time()
    print "pretime: ", pre_time
    CUR.execute(translatedQuery)
    post_time = datetime.datetime.now().time()
    print "post_time: ", post_time
    # CUR.fetchall()

    print "fetched rows:", CUR.row_count

    # for row in CUR.fetchall():
    #     print row