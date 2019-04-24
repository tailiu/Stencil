import proc_data as pd
import database as db
import graph as g

def lengthVsSize(stencilConnection, stencilCursor, destApp, startTime = None, endTime = None):
    if startTime == None or endTime == None:
        migrationIDs = pd.getAllMigrationIDs(stencilConnection, stencilCursor)
    else:
        migrationIDs = pd.getMigrationIDsBetweenTimestamps(stencilConnection, stencilCursor, startTime, endTime)
    length = []
    size = []
    for migrationID in migrationIDs:
        print migrationID
        l = pd.getLengthOfMigration(migrationID, stencilCursor)
        if  l == None:
            continue
        else:
            length.append(l)
            size.append(pd.getMigratedDataSize(destApp, migrationID, stencilCursor))
    g.line(size, length, "Migration Size (MB)", "Migration Length (s)")

def allLeftDataCumulativeGraph(stencilConnection, stencilCursor, destApp, srcApp):
    migrationIDs = pd.getAllMigrationIDs(stencilConnection, stencilCursor)
    l = []
    for migrationID in migrationIDs:
        print migrationID
        leftData1 = pd.getSizeOfLeftDataInMigratedRows(srcApp, migrationID, stencilCursor)
        leftData2 = pd.getSizeOfDataWithEntireRowLeft(srcApp, migrationID, stencilCursor)
        migratedData = pd.getMigratedDataSize(destApp, migrationID, stencilCursor)
        if leftData1 == None or leftData2 == None or migratedData == None:
            continue
        l.append((leftData1 + leftData2)/ float(migratedData + leftData1 + leftData2))
    print l
    g.cumulativeGraph(l, "Percentage of Data Left", "Probability")

stencilDB, srcApp, destApp, migrationID = "stencil", "diaspora", "mastodon", 1017008071
stencilConnection, stencilCursor = db.connDB(stencilDB)
# lengthVsSize(stencilConnection, stencilCursor, destApp)
# lengthVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 11:32:00", "2019-04-24 12:09:00")
allLeftDataCumulativeGraph(stencilConnection, stencilCursor, destApp, srcApp)

# print len(pd.getMigrationIDsBetweenTimestamps(stencilConnection, stencilCursor, "2019-04-24 11:32:00", "2019-04-24 12:09:00"))
# l = {0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.1, 0.1, 0.1, 0.1, 0.2, 0.2, 0.8}
# g.cumulativeGraph(l, "Percentage of Data Left", "Probability")
# print pd.getSizeOfLeftDataInMigratedRows(srcApp, migrationID, stencilCursor)
# print pd.getSizeOfDataWithEntireRowLeft(srcApp, migrationID, stencilCursor)
# print getAllMigrationIDs(stencilConnection, stencilCursor)
# print getMigratedDataSize(destApp, migrationID, stencilCursor)
# print getLengthOfMigration(migrationID, stencilCursor)
db.closeDB(stencilConnection)

