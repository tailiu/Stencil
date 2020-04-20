import psycopg2, copy
import numpy as np
from numpy import linalg as lg
import pandas as pd
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
from psycopg2.extras import RealDictCursor
from sklearn.metrics.pairwise import pairwise_distances, euclidean_distances
from sklearn.cluster import KMeans
from scipy.spatial.distance import pdist
import matplotlib.pyplot as plt

def getDBConn(db, cursor_dict=False):
    conn = psycopg2.connect(dbname=db, user="cow", password="123456", host="10.230.12.86", port="5432")
    conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT) 
    if cursor_dict is True:
        cursor = conn.cursor(cursor_factory = RealDictCursor)
    else:
        cursor = conn.cursor()
    return conn, cursor

def getPhysicalTables():
    tables = []
    tableq = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
    cur.execute(tableq)
    for row in cur.fetchall():
        table = row["tablename"]
        if table != "supplementary_tables" and ("supplementary_" in table or "base_" in table):
            tables.append(table)
    return tables

def createIndices():
    for table in getPhysicalTables():
        colsql = "select column_name from INFORMATION_SCHEMA.COLUMNS where table_name = '%s'"%table
        cur.execute(colsql)
        for row in cur.fetchall():
            column_name = row["column_name"]
            if column_name != "app_id" and (column_name == "pk" or column_name == "id" or "_id" in column_name):
                idxsql = "CREATE INDEX %s_%s_idx ON public.%s (%s); "%(table, column_name, table, column_name)
                # print idxsql
                cur.execute(idxsql)

def refreshMaterializedViews():

    sql = "select schemaname as schema_name, matviewname as view_name from pg_matviews;"
    cur.execute(sql)
    for row in cur.fetchall():
        sql = "REFRESH MATERIALIZED VIEW %s.%s WITH DATA;"%(row["schema_name"], row["view_name"])
        # print sql
        cur.execute(sql)

def deletePhysicalTables():
    for table in getPhysicalTables():
        q = 'DROP TABLE "%s" CASCADE;'%table
        cur.execute(q)

def truncatePhysicalTables():
    tables = ["supplementary_tables", "physical_schema", "physical_mappings", "migration_table"]
    for table in tables:
        # if table not in ["migration_registration", "display_flags", "txn_logs", "evaluation"]:
        sql = 'TRUNCATE "%s" RESTART IDENTITY CASCADE;'%table
        cur.execute(sql)

def getAppNameById(app_id):
    sql = "SELECT app_name FROM apps WHERE pk = " + str(app_id)
    cur.execute(sql)
    try: return cur.fetchone()["app_name"]
    except: return None

def getAppSchemas():
    sql = "SELECT app_id, table_name, column_name, app_schemas.pk as column_id, app_name FROM app_schemas JOIN app_tables ON app_tables.pk = app_schemas.table_id JOIN apps ON app_tables.app_id = apps.pk"
    cur.execute(sql)
    return cur.fetchall()

def getSchemaMappings():
    # sql = """
    #     SELECT as1.app_id AS app1, as1.table_name AS table1, as1.column_name AS col1, sm.source_attribute,
    #         as2.app_id AS app2, as2.table_name AS table2, as2.column_name AS col2, sm.dest_attribute
    #     FROM app_schemas as1 
    #     JOIN schema_mappings sm ON as1.pk = sm.source_attribute
    #     JOIN schema_mappings sm2 ON sm.source_attribute = sm2.dest_attribute AND sm.dest_attribute = sm2.source_attribute
    #     JOIN app_schemas as2 ON as2.pk = sm.dest_attribute
    #     WHERE sm.pk < sm2.pk 
    #     ORDER BY as1.app_id
    # """

    sql = """
            select a1.pk AS app1, at1.table_name AS table1, as1.column_name AS col1, sm.source_attribute, a2.pk AS app2, at2.table_name AS table2, as2.column_name AS col2, sm.dest_attribute
            from schema_mappings sm 
            join app_schemas as1 on sm.source_attribute = as1.pk 
            join app_tables at1 on as1.table_id = at1.pk
            join apps a1 on at1.app_id = a1.pk
            join app_schemas as2 on sm.dest_attribute = as2.pk
            join app_tables at2 on as2.table_id = at2.pk
            join apps a2 on at2.app_id = a2.pk
            ORDER BY a1.pk
    """

    cur.execute(sql)
    return cur.fetchall()

def genTransitivelyMappedAttrs(schema_mappings):

    common_attrs = {}

    for row in schema_mappings:
        srcAttr, dstAttr = row["source_attribute"], row["dest_attribute"]
        alreadyMapped = False
        for attr, mapped_attrs in common_attrs.items():
            if srcAttr in mapped_attrs:
                alreadyMapped = True
                break
            if dstAttr in mapped_attrs:
                common_attrs[attr].append(srcAttr)
                alreadyMapped = True
                break
        if alreadyMapped: continue
        # if row["dest_attribute"] in common_attrs.keys(): continue

        if srcAttr in common_attrs.keys():
            common_attrs[srcAttr].append(dstAttr)
        else:
            common_attrs[srcAttr] = [dstAttr]
        if dstAttr in common_attrs.keys():
            common_attrs[srcAttr] += common_attrs[dstAttr]
            common_attrs.pop(dstAttr)
                

        # for attr, mapped_attrs in common_attrs.items():
        #     if srcAttr in mapped_attrs or dstAttr in mapped_attrs:
        #         common_attrs[attr].append(dstAttr)
        #         common_attrs[attr].append(srcAttr)
        #         common_attrs[attr] = list(set(common_attrs[attr]))
        #         break
        # else:
        #     common_attrs[srcAttr] = [dstAttr]
        
        # if srcAttr in common_attrs.keys():
        #     common_attrs[srcAttr].append(dstAttr)
        # else:
        #     for attr, mapped_attrs in common_attrs.items():
        #         if srcAttr in mapped_attrs and dstAttr not in mapped_attrs:
        #             common_attrs[attr].append(dstAttr)
        #             break
        #     else:
        #         common_attrs[srcAttr] = [dstAttr]
    return common_attrs

def getAppSchemaAttr(attr):
    sql = '''
        SELECT app_id, table_name, column_name, app_schemas.pk as column_id, app_name 
        FROM app_schemas JOIN app_tables ON app_tables.pk = app_schemas.table_id JOIN apps ON app_tables.app_id = apps.pk
        WHERE app_schemas.pk = '%s'
    '''
    cur.execute(sql % attr)
    res = cur.fetchone()
    attr = "%s.%s.%s" % (res["app_name"], res["table_name"], res["column_name"])
    return attr

def printTransAttrs(trans_attrs):
    for attr, mapped_attrs in trans_attrs.items():
        # print attr, "=>", mapped_attrs
        print getAppSchemaAttr(attr), "=>",
        for mapped_attr in mapped_attrs:
            print getAppSchemaAttr(mapped_attr),
        print ""
    exit(0)

def filterSchemaRow(app_schemas, attr):
    return filter(lambda x: x.get('column_id') == attr, app_schemas)[0]

def genTableAppMapping(app_schemas, attrs):

    tables = {}

    for attr, mapped_attrs in attrs.items():
        row = filterSchemaRow(app_schemas, attr)
        app_id = row["app_id"]
        # app_name = row["app_name"]
        table_name = row["table_name"]
        col_name = row["column_id"]
        mapped_apps = [app_id]
        for mattr in mapped_attrs:
            mrow = filterSchemaRow(app_schemas,mattr)
            mapped_apps.append(mrow["app_id"])
        
        if table_name in tables.keys():
            if col_name in tables[table_name].keys():
                tables[table_name][col_name] += mapped_apps
            else:
                tables[table_name][col_name] = mapped_apps
        else:
            tables[table_name] = {col_name : mapped_apps}
    return tables

def genAttributeNodeVectors(app_schemas, attrs):
    
    tables = genTableAppMapping(app_schemas, attrs)
    node_vectors = {}
    for table_name, col_dict in tables.items():
        apps_using_table = list(reduce(lambda res, key: res.union(col_dict[key]), col_dict, set()))
        table_matrix, cols_inorder = [], []
        for col_name, col_apps in col_dict.items():
            cols_inorder.append(col_name)
            row_vector = [0] * len(apps_using_table)
            for col_app in col_apps:
                app_index = apps_using_table.index(col_app)
                if app_index >= 0:
                    row_vector[app_index] = 1
            table_matrix.append(row_vector)
        node_vectors[table_name] = pd.DataFrame( data = table_matrix, index = cols_inorder, columns = apps_using_table )
    return node_vectors

def genSimilarityMatrix(node_vector):
    return 1 - pairwise_distances(node_vector, metric = "cosine")

def genDistanceMatrix(node_vector):
    return euclidean_distances(node_vector)

def bendTheKnee(values):
    try:
        sd = [values[i+1] + values[i-1] - 2 * values[i] for i in range(1, len(values)-1)]
        return sd.index([max(sd)])
    except ValueError as e:
        print "Exception:", e, "| Returning: 1"
        return 1

def getNumberOfK(node_vector):

    dmat = euclidean_distances(node_vector)
    median = np.median(dmat)
    if median == 0: return 1
    gamma = (-1) / median
    kmat = np.zeros(dmat.shape)
    for x in range(0, dmat.shape[0]):
        for y in range(0, dmat.shape[0]):
            kmat[x,y] = np.exp(gamma * (dmat[x,y]**2))
    eigval, eigvec = lg.eig(kmat)
    v1N = np.full((dmat.shape[0], 1), 1.0/dmat.shape[0])
    
    vals = []
    indices = range(0, eigvec.shape[0]) 
    for i in indices:
        eVal = eigval[i] * np.square(np.dot(v1N.T, eigvec[i,:]))
        if eVal > 0:
            vals.append( np.log(eigval[i] * np.square(np.dot(v1N.T, eigvec[i,:]))))

    realValues = [abs(v.real) for v in vals]

    return bendTheKnee(realValues) # starts from 0, need to add 1

def genBaseTables(filtered_vector):
    
    optimal_clusters = getNumberOfK(filtered_vector) or 1
    kmeans = KMeans(n_clusters=optimal_clusters, random_state=0).fit(filtered_vector)
    base_tables = {}
    for i, k in enumerate(kmeans.labels_):
        k += 1
        if k in base_tables.keys():
            base_tables[k].append(filtered_vector.iloc[i].name)
        else:
            base_tables[k] = [filtered_vector.iloc[i].name]
    return base_tables

def createTable(name, attrs):

    isql =  "INSERT INTO PHYSICAL_SCHEMA (table_name, column_name) VALUES " + \
            ", ".join(["('%s','%s')" % (name, attr.split()[0]) for attr in attrs])

    attrs.insert(0,"pk SERIAL PRIMARY KEY")
    attrs.append("base_mark_delete BOOL")
    tsql = "CREATE TABLE %s ( %s )" % (name, ', '.join(attrs))
    
    # print tsql, isql

    # cur.execute(tsql)
    # cur.execute(isql)

def createBaseTable(name, attrs, app_schemas, trans_attrs):

    attrs = {a: filterSchemaRow(app_schemas, a)["column_name"] for a in attrs}

    # isql =  "INSERT INTO PHYSICAL_SCHEMA (table_name, column_name) VALUES " + ", ".join(["('%s','%s')" % (name, attr) for attr in attrs])

    for attr_id, attr_name in attrs.items():
        isql =  "INSERT INTO PHYSICAL_SCHEMA (table_name, column_name) VALUES ('%s', '%s') RETURNING pk" % (name, attr_name)
        # print isql
        cur.execute(isql)
        phy_attr_id = cur.fetchone()['pk']
        mapped_attrs = [attr_id] + trans_attrs[attr_id]
        pmsql = "INSERT INTO PHYSICAL_MAPPINGS (logical_attribute, physical_attribute) VALUES " + ", ".join(["('%s', '%s')" % (attr_id,phy_attr_id) for attr_id in set(mapped_attrs)])
        # print pmsql
        cur.execute(pmsql)

    attrs_with_type = [attr + " varchar" for attr in attrs.values()]
    # attrs_with_type.insert(0,"app_id varchar")
    # attrs_with_type.insert(0,"base_pk SERIAL PRIMARY KEY")
    attrs_with_type.insert(0,"pk int8 PRIMARY KEY")
    # attrs_with_type.append("base_created_at TIMESTAMP DEFAULT now()")
    # attrs_with_type.append("base_mark_delete BOOL")
    tsql = "CREATE TABLE %s ( %s )" % (name, ', '.join(attrs_with_type))
    # print tsql
    cur.execute(tsql)

def createSupplementaryTables():

    sql = """
        SELECT  app_id, table_id, string_agg(concat('"',column_name, '"', ' ', data_type),',') as column_names, string_agg(app_schemas.pk::varchar,',') as column_ids
        FROM app_schemas JOIN app_tables ON app_schemas.table_id = app_tables.pk 
        WHERE app_schemas.pk NOT IN (SELECT logical_attribute FROM physical_mappings) 
        GROUP BY app_id, table_id
        ORDER BY app_id
    """
    cur.execute(sql)

    for row in cur.fetchall():

        app_id      = row["app_id"]
        table_id    = row["table_id"]
        column_names= row["column_names"]
        column_ids  = row["column_ids"]
        # print app_id, table_id, column_names

        insql = "INSERT INTO SUPPLEMENTARY_TABLES (table_id) VALUES(%d) RETURNING pk" % (table_id)
        # print insql
        cur.execute(insql)
        
        supp_table_id = cur.fetchone()['pk']
        
        cols = [attr  for attr in column_names.split(',') if len(attr) and attr != "app_id"]
        cols.insert(0,"pk int8 PRIMARY KEY")
        # cols.append("supp_created_at TIMESTAMP DEFAULT now()")
        # cols.append("supp_mark_delete BOOL")
        
        tsql = "CREATE TABLE %s ( %s )" % ("supplementary_%s"%supp_table_id, ', '.join(cols))
        # print tsql
        cur.execute(tsql)

if __name__ == "__main__":

    t = 0.5

    db, cur = getDBConn("stencil_test", True)
    
    print "Reset Physical DB"
    truncatePhysicalTables()
    deletePhysicalTables()

    print "Get App Schemas"
    app_schemas = getAppSchemas()

    print "Get Schema Mappings"
    schema_mappings = getSchemaMappings()

    print "Transitively Mapped Attrs"
    trans_attrs = genTransitivelyMappedAttrs(schema_mappings)

    # print printTransAttrs(trans_attrs)

    print "Attribute Node Vectors"
    node_vectors = genAttributeNodeVectors(app_schemas, trans_attrs)

    print "createBaseTables"
    for table, vector in node_vectors.items():
        filtered_vector = vector.loc[vector.sum(axis=1)/vector.shape[1] >= t]            
        base_tables = genBaseTables(filtered_vector)
        for idx, base_attrs in base_tables.items():
            bt_name = "base_%s_%s" % (table, idx)
            # print bt_name
            createBaseTable(bt_name, base_attrs, app_schemas, trans_attrs)

    print "createSupplementaryTables"    
    createSupplementaryTables()

    print "createIndices"    
    createIndices()

    print "refreshMaterializedViews"
    refreshMaterializedViews()