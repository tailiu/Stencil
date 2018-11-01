import re
import utils
import datetime
from db import DB


def resolveRequest(query, baseAttributes, suppAttributes):
    for attr in baseAttributes: query = re.sub(r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[2].lower(), query)
    for attr in suppAttributes: query = re.sub(r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[0].lower(), query)

    tables = set()
    for attr in baseAttributes: tables.add(attr[1])
    for attr in suppAttributes: tables.add(attr[1])

    if len(tables) == 1: fromTables = tables.pop()
    else:
        fromTables = ''
        table1 = tables.pop()
        fromTables += table1
        while len(tables) > 0:
            table2 = tables.pop()
            fromTables += ' join ' + table2 + ' on ' + table1 + '.row_id = ' + table2 + '.row_id '
            if len(tables) > 0: table1 = table2

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

    db = DB()
    CONN, CUR = db.conn, db.cursor

    app = 'hacker news'

    sql = "SELECT By, Descendents, Id, Retrieved_on, Score\
           FROM story  \
           WHERE By = 'Impossible' and Id = 13075839" 
    
    sql1 = "SELECT * \
            FROM story  \
            WHERE By = 'Impossible' and Id = 13075839"

    sql2 = 'SELECT * \
            FROM comment  \
<<<<<<< HEAD
            WHERE user = "lisper" '
=======
            WHERE By = 'lisper' "
>>>>>>> d819b464149ab081a8313187c5b69ee3f129a717

    sql3 = "SELECT * \
            FROM tweet"

    sql4 = "SELECT * from story"

    sql4 = "SELECT * \
            FROM comment"

    sql5 = "SELECT * \
            FROM story  \
            WHERE By = 'lisper' "

    translatedQuery = translateBasicSelectQuery(CUR, sql2, app)
    print "translatedQuery: ", translatedQuery

    pre_time = datetime.datetime.now()
    CUR.execute(translatedQuery)
    post_time = datetime.datetime.now()
    print "pretime: %s; post time: %s" % (pre_time, post_time)
    print "Time duration: {0}".format(post_time - pre_time) 

    print "fetched rows:", len(CUR.fetchall())
    # for row in CUR.fetchall():
    #     print row

    db.close()