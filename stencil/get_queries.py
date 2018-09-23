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
    if query.find('where') == -1: query = query[:tablesStart] + ' ' + fromTables
    else: query = query[:tablesStart] +  ' ' + fromTables + ' ' + query[query.find("where"):]

    return query

def translateBasicSelectQuery(CUR, query, app):
    query = query.lower()

    condList = []
    if query.find('where') == -1:
        table = utils.findBetweenStrings(query, 'from', None).strip() # Assume there is only one table
    else:
        table = utils.findBetweenStrings(query, 'from', 'where').strip() # Assume there is only one table
        condList = utils.processConditions(utils.findBetweenStrings(query, 'where', None))
        condList = utils.removeSpace(condList)

    print table

    if query.find('*') == -1:
        attributes = utils.findBetweenStrings(query, 'select', 'from').split(',')
        attributes = utils.removeSpace(attributes)
    else: 
        attributes = utils.findAllAttributes(CUR, app, table)
        attrStr = ''
        for attr in attributes: attrStr += ", " + attr
        attrStr = attrStr.replace(',', '', 1)
        query = query.replace('*', attrStr)
        query = query.lower()

    if len(condList): attributes = list(set(attributes).union(condList))

    baseAttributes = utils.translateAttributesToBaseTables(CUR, app, table, attributes)

    suppAttributeList = []
    for attr in attributes:
        find = False
        for var in baseAttributes:
            if attr.lower() == var[0].lower(): 
                find = True
                break
        if not find: suppAttributeList.append(attr)

    suppAttributes = ()
    if len(suppAttributeList) != 0: suppAttributes = utils.findSuppTables(CUR, app, table, suppAttributeList)

    return resolveRequest(query, baseAttributes, suppAttributes)


if __name__ == "__main__":

    CONN, CUR = utils.getDBConn()

    app = 'hacker news'

    sql = "SELECT By, Descendents, Id, Retrieved_on, Score\
           FROM story  \
           WHERE By = 'Impossible' and Id = 13075839" 
    
    sql1 = "SELECT * \
            FROM story  \
            WHERE By = 'Impossible' and Id = 13075839"

    sql2 = "SELECT * \
            FROM comment  \
            WHERE user = 'lisper' "

    sql3 = "SELECT * \
            FROM comment"

<<<<<<< HEAD
    translatedQuery = translateBasicSelectQuery(CUR, sql3)
    print translatedQuery

    translatedQuery = translateBasicSelectQuery(CUR, sql3)
=======
    translatedQuery = translateBasicSelectQuery(CUR, sql3, app)
>>>>>>> f1b841de68d08b02c71cafc1bc3f9abef71bfafd
    print "translatedQuery: ", translatedQuery

    pre_time = datetime.datetime.now().time()
    print "pretime: ", pre_time
    CUR.execute(translatedQuery)
    post_time = datetime.datetime.now().time()
    print "pretime: %s; post time: %s" % (pre_time, post_time)

    print "fetched rows:", len(CUR.fetchall())

    # for row in CUR.fetchall():
    #     print row