import MySQLdb
import utils
import re

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

def resolveRequest(query, baseAttributes, suppAttributes):
    print baseAttributes
    print suppAttributes
    for attr in baseAttributes: query = re.sub("(?<=\s|=)"r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[2].lower(), query)
    for attr in suppAttributes: query = re.sub("(?<=\s|=)"r"\b{0}\b".format(attr[0].lower()), attr[1].lower() + '.' + attr[0].lower(), query)

    tables = set()
    for attr in baseAttributes: tables.add(attr[1])
    for attr in suppAttributes: tables.add(attr[1])

    if len(tables) == 1:
        fromTables = tables.pop()
        oneTable = fromTables
    else:
        fromTables = ''
        for table in tables: 
            fromTables += ' inner join ' + table
            oneTable = table
        fromTables = fromTables.replace('inner join', '', 1)
        fromTables += ' on '
        for table in tables: fromTables += ' = ' + table + '.row_id '
        fromTables = fromTables.replace('=', '', 1)

    tablesStart = query.find("from") + len("from")
    tablesEnd = query.find("where")
    query = query[:tablesStart] +  ' ' + fromTables + ' ' + query[tablesEnd:]

    return query

def processJoinConditions(query):
    conds = utils.findBetweenStrings(query, 'on', 'where')
    conds = re.split('and | or', conds)

    condList = []
    for cond in conds:
        cond1 = cond.split('=')
        condList.append(cond1[0].strip())
        condList.append(cond1[1].strip())
    return condList

def changeQuery(query):
    conds = utils.findBetweenStrings(query, 'on', 'where')
    
    if query.find('where') == -1: return query + ' where ' + conds
    else: return query + ' and ' + conds 

def translateJoinQuery(CUR, query):
    query = query.lower()

    tableInfo = utils.findBetweenStrings(query, 'from', 'where').strip()
    tables = tableInfo[:tableInfo.find('on')].split('inner join')
    tables = utils.removeSpace(tables)
    
    if query.find('*') == -1:
        attributes = utils.findBetweenStrings(query, 'select', 'from').split(',')
        attributes = utils.removeSpace(attributes)
    # else: 
    #     attributes = utils.findAllAttributes('hacker news', tables)
    #     attrStr = ''
    #     for attr in attributes: attrStr += ", " + attr
    #     attrStr = attrStr.replace(',', '', 1)
    #     query = query.replace('*', attrStr)
    #     query = query.lower()

    condList1 = utils.processConditions(utils.findBetweenStrings(query, 'where', None))
    condList1 = utils.removeSpace(condList1)
    condList2 = processJoinConditions(query)
    query = changeQuery(query)

    attrList = list(set(attributes).union(condList1))
    attrList = list(set(attrList).union(condList2))

    baseAttributes = utils.translateAttributesToBaseTables(CUR, 'hacker news', tables, attrList)

    suppAttributeList = []
    for attr in attrList:
        find = False
        for var in baseAttributes:
            if attr.lower() == var[0].lower(): 
                find = True
                break
        if not find: suppAttributeList.append(attr)

    suppAttributes = ()
    if len(suppAttributeList) != 0: suppAttributes = utils.findSuppTables(CUR, 'hacker news', tables, suppAttributeList)

    return resolveRequest(query, baseAttributes, suppAttributes)


if __name__ == "__main__":

    CONN, CUR = utils.getDBConn()

    sql = "SELECT descendents, kids, parent, story.id\
            FROM story INNER JOIN comment on story.id = comment.parent \
            WHERE story.By = 'lisper'"
    sql1 = "SELECT descendents, kids, parent, story.id\
            FROM story INNER JOIN comment on story.id = comment.parent \
            WHERE story.By = 'lisper'"

    sql2 = "select base_1.parent, base_1.id\
            from base_1 inner join base_1 as B on base_1.id = B.parent\
            where base_1.user = 'lisper'"

    translatedQuery = translateJoinQuery(CUR, sql1)

    print translatedQuery
    CUR.execute(sql2)

    for row in CUR.fetchall():
        print row