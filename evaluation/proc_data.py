import matplotlib.pyplot as plt
import numpy as np
import database as db

def _getAColumnSize(cursor, cols, table, id):
    select = ""
    for i, v in enumerate(cols):
        select += "COALESCE(pg_column_size({}), 0)".format(v)
        if i != len(cols) - 1:
            select += "+"
    query = "select {} from {} where id = {}".format(select, table, id)
    try:
        result = db.getDataFromDatabase(cursor, query)[0][0]
        return result
    except:
        print "Query Error: {}".format(query)

def _getAllDataBasedOnMigrationID(migrationID, stencilCursor):
    query = "select * from evaluation where migration_id = '{}'".format(migrationID)
    return db.getDataFromDatabase(stencilCursor, query)

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
    totalSize = 0
    destAppConn, destAppCursor = db.connDB(destApp)
    migratedData = _getAllDataBasedOnMigrationID(migrationID, stencilCursor)
    for row in migratedData:
        if row[8] == "n/a":
            continue
        size = _getAColumnSize(destAppCursor, row[8].split(","), row[4], row[6])
        totalSize += float(size) / 10**3
    db.closeDB(destAppConn)
    return totalSize

def getAllMigrationIDs(stencilConnection, stencilCursor):
    query = "select distinct on (migration_id) migration_id from evaluation"
    rows = db.getDataFromDatabase(stencilCursor, query)
    migrationIDs = []
    for row in rows:
        migrationIDs.append(row[0])
    return migrationIDs

def getMigrationIDsBetweenTimestamps(stencilConnection, stencilCursor, startTime, endTime):
    query = "select distinct on (action_id) action_id from txn_logs where created_at >= timestamp '{}' and created_at <= timestamp '{}';".format(startTime, endTime)
    rows = db.getDataFromDatabase(stencilCursor, query)
    migrationIDs = []
    for row in rows:
        migrationIDs.append(row[0])
    return migrationIDs

def getSizeOfLeftDataInMigratedRows(srcApp, migrationID, stencilCursor):
    migratedData = _getAllDataBasedOnMigrationID(migrationID, stencilCursor)
    totalSize = 0
    srcAppConn, srcAppCursor = db.connDB(srcApp)
    leftData = {}

    for row in migratedData:
        if row[8] == "n/a":
            continue
        table = row[3]
        key = table + ":" + row[5]
        migratedCols = set(row[7].split(","))
        allCols = set(db.getColsOfTable(srcAppCursor, table))
        leftCols = allCols - migratedCols
        if key in leftData:
            leftData[key] = leftCols.intersection(leftData[key])
        else:
            leftData[key] = leftCols
    
    for key in leftData:
        size = _getAColumnSize(srcAppCursor, leftData[key], key.split(":")[0], key.split(":")[1])
        if size == None:
            return
        totalSize += float(size) / 10**3
    
    db.closeDB(srcAppConn)
    return totalSize

def getSizeOfDataWithEntireRowLeft(srcApp, migrationID, stencilCursor):
    migratedData = _getAllDataBasedOnMigrationID(migrationID, stencilCursor)
    srcAppConn, srcAppCursor = db.connDB(srcApp)
    totalSize = 0

    for row in migratedData:
        if row[8] != "n/a":
            continue
        table = row[3]
        allCols = db.getColsOfTable(srcAppCursor, table)
        size = _getAColumnSize(srcAppCursor, allCols, row[3], row[5])
        if size == None:
            return
        totalSize += float(size) / 10**3
    
    db.closeDB(srcAppConn)
    return totalSize

def getPercentageInIntervals(data, step):
    totalSize = len(data)
    percentage = []
    currStep = step
    while currStep <= 1.00000001:
        num = 0
        for d in data:
            if currStep == 1.0:
                if d <= currStep and d >= currStep - step:
                    num += 1
            else:
                if d < currStep and d >= currStep - step:
                    num += 1
        percentage.append(num/float(totalSize))
        currStep += step
    return percentage
    