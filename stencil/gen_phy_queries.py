import MySQLdb, json, sqlparse, uuid
from gen_input_queries import CONN, CUR, getSchemaMapping, getTableNames

def getNewRowId():
    return uuid.uuid4().hex

def getAppId(app_name):
    sql = "SELECT PK from apps WHERE app_name = '%s'" % app_name
    CUR.execute(sql)
    return CUR.fetchone()[0]

def getColValDict(tokens):
    cols, vals = [], []
    for i in range(3, len(tokens)): # column names start from index 3
        if tokens[i].lower() == "values":  break
        cols.append(tokens[i])
    for i in range(4+len(cols), len(tokens)): # values start from 3 + number of columns + 1 (VALUES)
        vals.append(tokens[i])
    return { x:y for x,y in zip(cols, vals)}

def getQueriesFromPQIng(pq_ing, app_name):
    row_id = getNewRowId()
    app_id = getAppId(app_name)
    sqls   = []
    for tn in pq_ing.keys():
        sql = 'INSERT INTO %s ( row_id, app_id, ' % tn
        val = ' VALUES ( "%s", %s, ' % (row_id, app_id)
        for tc in pq_ing[tn].keys():
            sql += '%s,' %tc
            val += '"%s",' % CONN.escape_string(pq_ing[tn][tc])
        sqls.append(sql[:-1] + ")" + val[:-1] + ");")
    return sqls

if __name__ == "__main__":
    
    queries = []
    with open("hn_log.queries") as fh: queries = fh.read().split('\n')[:-1]
    
    app_name = "hacker news"
    schemas  = {tn : getSchemaMapping(app_name, tn) for tn in getTableNames(app_name)}

    pqueries = []

    for q in queries:
        
        tokens = [x.value for x in list(sqlparse.parse(q)[0].flatten()) if x.value.strip(". (),") != ""]
        tname  = tokens[2]
        colval = getColValDict(tokens)
        
        pq_ing = {}

        if tname in schemas.keys():
            for col in colval.keys():
                if col in schemas[tname].keys():
                    phy_tn  = schemas[tname][col][0]
                    phy_col = schemas[tname][col][1]
                    if phy_tn in pq_ing.keys():
                        pq_ing[phy_tn][phy_col] = colval[col]
                    else:
                        pq_ing[phy_tn] = {phy_col : colval[col]}

        pqs = getQueriesFromPQIng(pq_ing, app_name)
        for pq in pqs:
            # print pq
            CUR.execute(pq)
        CONN.commit()
