import MySQLdb
import sqlparse
import uuid
import json
import re
from db import DB

class QueryResolver():

    def __init__(self, app_name):
        self.db       = DB()
        self.app_name = app_name
        self.app_id   = self.__getAppId(app_name)
        self.base_map = self.__getBaseMappings()
        self.supl_map = self.__getSupplementaryMappings()
        self.pqs      = []

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
        return self.db.cursor.fetchall()
        return {(row[0], row[1]) : (row[2], row[3]) for row in rows}
    
    def __getSupplementaryMappings(self):
        sql = """   SELECT  LOWER(app_tables.table_name), 
                            LOWER(app_schemas.column_name), 
                            LOWER(supplementary_tables.supplementary_table),
                            LOWER(app_schemas.column_name)
                    FROM 	app_tables JOIN 
                    		app_schemas ON app_schemas.table_id = app_tables.PK JOIN
                    		supplementary_tables ON supplementary_tables.table_id = app_tables.PK
                    WHERE 	app_tables.app_id  = "%s" AND
                    		app_schemas.PK NOT IN (
                                SELECT logical_attribute FROM physical_mappings
                            )""" % (self.app_id) 
        self.db.cursor.execute(sql)
        return self.db.cursor.fetchall()
        return {(row[0], row[1]) : (row[2], row[1]) for row in rows}

    def __getInsertQueryIngs(self, q):
        
        tokens   = [x.value for x in list(sqlparse.parse(q)[0].flatten()) if x.value.strip(". (),") != ""]
        of_tname = tokens.index("INTO")+1
        of_cols  = of_tname + 1
        of_vals  = tokens.index("VALUES")+1
        tname    = tokens[of_tname].lower()
        ings     = { "table": tname , "items": {} }

        for i, j in zip(range(of_cols, of_vals), range(of_vals, len(tokens))):
            ings["items"][tokens[i].lower()] = tokens[j]

        return ings

    def __getUpdateQueryIngs(self, q):
        
        tokens = re.split('(update table | set | where)(?i)', q)
        # todo
        return tokens
    
    def __getDeleteQueryIngs(self, q):
        
        regex  = "(delete from | where )(?i)"
        tokens = filter(lambda x: x.strip(" ,"), re.split(regex, q))
        table  = tokens[1]
        conds  = tokens[3]
        phy_tabs = self.__getPhyMappingForLogicalTable(table)

        print phy_tabs
        # print table
        # print conds
        # todo
        # return tokens
    
    def __getPhyMappingForLogicalTable(self, ltable):

        phy_map = {}

        for item in self.base_map + self.supl_map:
            if item[0].lower() == ltable:
                base_tab = item[2].lower()
                base_col = item[3].lower()
                log_col  = item[1].lower()
                if base_tab in phy_map.keys():
                    phy_map[base_tab].append((base_col , log_col))
                else:
                    phy_map[base_tab] = [(base_col , log_col)]
        
        return phy_map

    def sendToDB(self):
        if self.pqs:
            for pq in self.pqs:
                print pq
                self.db.cursor.execute(pq)
            self.db.conn.commit()
            del self.pqs[:]

    def resolveInsert(self, q):
        
        row_id  = self.__getNewRowId()
        q_ing   = self.__getInsertQueryIngs(q)
        phy_map = self.__getPhyMappingForLogicalTable(q_ing["table"])

        for pt in phy_map.keys():
            pq_cols = 'INSERT INTO %s ( app_id, Row_id,' % pt
            pq_vals = ' VALUES ( %s, "%s",' % (self.app_id, row_id)
            for phy_col in phy_map[pt]:
                if phy_col[1] in q_ing["items"].keys():
                    pq_cols += '%s,' % phy_col[0]
                    pq_vals += '%s,' % q_ing["items"][phy_col[1]]
            pq = pq_cols[:-1] + ")" + pq_vals[:-1] + ");"
            self.pqs.append(pq) 
        return self.pqs

    def resolveUpdate(self, q):
        
        q_ing    = self.__getUpdateQueryIngs(q)
        print q_ing

    def resolveDelete(self,q):
        
        self.__getDeleteQueryIngs(q)