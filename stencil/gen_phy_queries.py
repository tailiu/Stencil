import MySQLdb, json, sqlparse, uuid
from db import DB

class InsertQueryResolver():

    def __init__(self, app_name, q):
        self.db       = DB()
        self.app_name = app_name
        self.app_id   = self.__getAppId(app_name)
        self.q        = q
        self.row_id   = self.__getNewRowId()

    def __del__(self):
        self.db.close()

    def __getNewRowId(self):
        return uuid.uuid4().hex

    def __getAppId(self, app_name):
        sql = "SELECT PK from apps WHERE app_name = '%s'" % app_name
        self.db.cursor.execute(sql)
        return self.db.cursor.fetchone()[0]

    def __getBaseMappings(self):
        sql = """ SELECT    LOWER(app_tables.table_name), 
                            LOWER(app_schemas.column_name), 
                            LOWER(base_table_attributes.table_name), 
                            LOWER(base_table_attributes.column_name)
                    FROM 	physical_mappings 
                    JOIN 	app_schemas ON physical_mappings.logical_attribute = app_schemas.PK
                    JOIN 	app_tables ON app_schemas.table_id = app_tables.PK
                    JOIN 	base_table_attributes ON physical_mappings.physical_attribute = base_table_attributes.PK
                    WHERE 	app_tables.app_id  = "%s" """ % (self.app_id) 
        self.db.cursor.execute(sql)
        rows = self.db.cursor.fetchall()
        return {(row[0], row[1]) : (row[2], row[3]) for row in rows}
    
    def __getSupplementaryMappings(self):
        sql = """   SELECT  LOWER(app_tables.table_name), 
                            LOWER(app_schemas.column_name), 
                            LOWER(supplementary_tables.supplementary_table)
                    FROM 	app_tables JOIN 
                    		app_schemas ON app_schemas.table_id = app_tables.PK JOIN
                    		supplementary_tables ON supplementary_tables.table_id = app_tables.PK
                    WHERE 	app_tables.app_id  = "%s" AND
                    		app_schemas.PK NOT IN (
                                SELECT logical_attribute FROM physical_mappings
                            )""" % (self.app_id) 
        self.db.cursor.execute(sql)
        rows = self.db.cursor.fetchall()
        return {(row[0], row[1]) : (row[2], row[1]) for row in rows}

    def __getQueryIngs(self):
        
        tokens   = [x.value for x in list(sqlparse.parse(self.q)[0].flatten()) if x.value.strip(". (),") != ""]
        of_tname = tokens.index("INTO")+1
        of_cols  = of_tname + 1
        of_vals  = tokens.index("VALUES")+1
        tname    = tokens[of_tname].lower()
        ings     = {}

        for i, j in zip(range(of_cols, of_vals), range(of_vals, len(tokens))):
            ings[(tname, tokens[i].lower())] = tokens[j]

        return ings

    def resolve(self, q, send_to_db = False):
        
        base_map = self.__getBaseMappings()
        supl_map = self.__getSupplementaryMappings()
        phy_map  = dict(base_map.items() + supl_map.items())
        phy_tabs = list(set([x[0] for x in phy_map.values()]))
        q_ing    = self.__getQueryIngs()
        pqs      = []

        for pt in phy_tabs:
            valid   = False
            pq_cols = 'INSERT INTO %s ( app_id, Row_id,' % pt
            pq_vals = ' VALUES ( %s, "%s",' % (self.app_id, self.row_id)
            for item in phy_map.items():
                if item[1][0] == pt:
                    if item[0] in q_ing.keys():
                        valid = True
                        pq_cols += '%s,' % item[1][1]
                        pq_vals += '%s,' % q_ing[item[0]]
            if valid:
                pq = pq_cols[:-1] + ")" + pq_vals[:-1] + ");"
                pqs.append(pq)

        if send_to_db:
            for pq in pqs:
                # print pq
                self.db.cursor.execute(pq)
            self.db.conn.commit()
        return pqs

if __name__ == "__main__":

    app_name = "hacker news"
    with open("hn_log.queries") as fh: 
        queries = fh.read().split('\n')
        for q in queries:
            IQR = InsertQueryResolver(app_name, q)
            IQR.resolve(q, True)
            # break
        
        
