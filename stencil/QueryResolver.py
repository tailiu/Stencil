import MySQLdb
import sqlparse
import uuid
import json
import re
from db import DB
from utils import getRowID

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
        
        tokens   = [x.value for x in list(sqlparse.parse(q.strip(";"))[0].flatten()) if x.value.strip(". (),") != ""]
        of_tname = tokens.index("INTO")+1
        of_cols  = of_tname + 1
        of_vals  = tokens.index("VALUES")+1
        tname    = tokens[of_tname].lower()
        ings     = { "table": tname.strip() , "items": {} }

        for i, j in zip(range(of_cols, of_vals), range(of_vals, len(tokens))):
            ings["items"][tokens[i].strip("\"'` ").lower()] = tokens[j]

        return ings

    def __getUpdateQueryIngs(self, q):
        
        regex  = "(update | set | where )(?i)"
        tokens = filter(lambda x: x.strip(" ,"), re.split(regex, q))
        updates_list = [token.strip(" ,") for token in tokens[3].split(",")]
        updates_dict = {update.split("=")[0].strip().lower():update.split("=")[1].strip() for update in updates_list}
        return {"table": tokens[1].strip(), "conditions": tokens[-1], "updates": updates_dict}
    
    def __getDeleteQueryIngs(self, q):
        
        regex  = "(delete from | where )(?i)"
        tokens = filter(lambda x: x.strip(" ,"), re.split(regex, q))
        return {"table": tokens[1].strip(), "conditions": tokens[3]}
    
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

    def __get_affected_row_ids(self, ltable, conds):
        row_ids = getRowID(self.db.cursor, self.app_name, ltable, conds)
        return [x[0] for x in row_ids]

    def DBCommit(self):
        self.db.conn.commit()

    def runQuery(self):
        if self.pqs:
            for pq in self.pqs:
                # print pq
                self.db.cursor.execute(pq)
            del self.pqs[:]

    def resolveInsert(self, q):
        
        row_id  = self.__getNewRowId()
        q_ing   = self.__getInsertQueryIngs(q)
        phy_map = self.__getPhyMappingForLogicalTable(q_ing["table"])

        # exit()
        # print q_ing["items"], "\n"

        for pt in phy_map.keys():
            pq_cols = 'INSERT INTO `%s` ( app_id, Row_id,' % pt
            pq_vals = ' VALUES ( %s, "%s",' % (self.app_id, row_id)
            for phy_col in phy_map[pt]:
                if phy_col[1] in q_ing["items"].keys():
                    pq_cols += '`%s`,' % phy_col[0]
                    pq_vals += '%s,' % q_ing["items"][phy_col[1]]
            pq = pq_cols[:-1] + ")" + pq_vals[:-1] + ");"
            self.pqs.append(pq) 
        return self.pqs

    def resolveUpdate(self, q):
        
        ings    = self.__getUpdateQueryIngs(q)
        phy_map = self.__getPhyMappingForLogicalTable(ings["table"])
        row_ids = self.__get_affected_row_ids(ings["table"], ings["conditions"])

        if len(row_ids) <= 0: return

        for pt in phy_map.keys():
            updates = ""
            for mapping in phy_map[pt]:
                if mapping[1] in ings["updates"].keys():
                    updates += "`%s` = %s," % (mapping[0], ings["updates"][mapping[1]])
            if updates != "":
                updates = updates.strip(",")
                pq = 'UPDATE `%s` SET %s WHERE row_id IN (%s);'% (pt, updates, str(row_ids).strip("[]"))
                self.pqs.append(pq)
                # print pq

    def resolveDelete(self,q):
        
        ings    = self.__getDeleteQueryIngs(q)
        phy_map = self.__getPhyMappingForLogicalTable(ings["table"])
        row_ids = self.__get_affected_row_ids(ings["table"], ings["conditions"])

        if len(row_ids) <= 0: return

        for pt in phy_map.keys():
            pq = 'DELETE FROM %s WHERE row_id IN (%s);' % (pt, str(row_ids).strip("[]"))
            self.pqs.append(pq)
    
    def getResolvedQueries(self):
        return self.pqs
        