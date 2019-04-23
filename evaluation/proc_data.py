import matplotlib.pyplot as plt
import numpy as np
import database as db

def getLengthOfMigration(migrationID, stencilCursor):
    try:
        query1 = "select created_at from txn_logs where action_id = {} and action_type = 'BEGIN_TRANSACTION'".format(migrationID)
        query2 = "select created_at from txn_logs where action_id = {} and action_type = 'COMMIT'".format(migrationID)
        startTime = db.getDataFromDatabase(stencilCursor, query1)[0][0]
        endTime = db.getDataFromDatabase(stencilCursor, query2)[0][0]
        return (endTime - startTime).total_seconds()
    except IndexError:
        print "Error: Get the length of a migration {}".format(migrationID)

def getMigratedDataSize(destApp, migrationID, stencilCursor):
    query1 = "select * from evaluation where migration_id = '{}'".format(migrationID)
    migratedData = db.getDataFromDatabase(stencilCursor, query1)
    totalSize = 0
    destAppConn, destAppCursor = db.connDB(destApp)
    for row in migratedData:
        cols = row[8].split(",")
        select = ""
        for col in cols:
            query2 = "select pg_column_size({}) from {} where id = {}".format(col, row[4], row[6])
            size = db.getDataFromDatabase(destAppCursor, query2)[0][0]
            if size != None:
                totalSize += size
    db.closeDB(destAppConn)
    return totalSize

def getAllMigrationIDs(stencilConnection, stencilCursor):
    query = "select distinct on (migration_id) migration_id from evaluation"
    rows = db.getDataFromDatabase(stencilCursor, query)
    migrationIDs = []
    for row in rows:
        migrationIDs.append(row[0])
    return migrationIDs
    
    