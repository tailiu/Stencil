import proc_data as pd
import database as db
import graph as g
import datetime

def _log(fileName, l):
    f = open(fileName, "a")
    f.write("************************************\n")
    f.write(str(datetime.datetime.now()) + "\n")
    for data in l:
        f.write("[")
        f.write(", ".join([str(i) for i in data]))
        f.write("]\n")
    f.write("************************************\n")
    f.close()

def timeVsSize(stencilConnection, stencilCursor, destApp, startTime = None, endTime = None):
    if startTime == None or endTime == None:
        migrationIDs = pd.getAllMigrationIDs(stencilConnection, stencilCursor)
    else:
        migrationIDs = pd.getMigrationIDsBetweenTimestamps(stencilConnection, stencilCursor, startTime, endTime)
    time = []
    size = []
    for migrationID in migrationIDs:
        print migrationID
        l = pd.getLengthOfMigration(migrationID, stencilCursor)
        if  l == None:
            continue
        else:
            time.append(l)
            size.append(pd.getMigratedDataSize(destApp, migrationID, stencilCursor))
    _log("timeVsSize", [time, size])
    g.line(size, time, "Migration Size (KB)", "Migration Time (s)", "Migration Time Vs Migration Size")


def leftDataCumulativeGraph(stencilConnection, stencilCursor, destApp, srcApp, dataLeftInBrokenRows = False):
    migrationIDs = pd.getAllMigrationIDs(stencilConnection, stencilCursor)
    l = []
    for migrationID in migrationIDs:
        print migrationID
        leftData1 = pd.getSizeOfLeftDataInMigratedRows(srcApp, migrationID, stencilCursor)
        leftData2 = 0
        if not dataLeftInBrokenRows:
            leftData2 = pd.getSizeOfDataWithEntireRowLeft(srcApp, migrationID, stencilCursor)
        migratedData = pd.getMigratedDataSize(destApp, migrationID, stencilCursor)
        if leftData1 == None or leftData2 == None or migratedData == None:
            continue
        l.append((leftData1 + leftData2)/ float(migratedData + leftData1 + leftData2))
    if not dataLeftInBrokenRows:
        _log("allLeftData", [l])
        g.cumulativeGraph(l, "Percentage of Data Left", "Probability")
    else:
        _log("dataLeftInBrokenRows", [l])
        g.cumulativeGraph(l, "Percentage of Data Left in Broken Rows", "Probability")

stencilDB, srcApp, destApp, migrationID = "stencil", "diaspora", "mastodon", 1017008071
stencilConnection, stencilCursor = db.connDB(stencilDB)
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-03-24 11:32:00", "2019-04-24 11:31:00") # 1 thread
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 11:32:00", "2019-04-24 12:09:00") # 5 threads
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 12:10:00", "2019-04-24 13:04:00") # 10 threads
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 13:09:00", "2019-04-24 15:45:00") # 20 threads
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-24 15:46:00", "2019-04-25 10:43:00") # 50 threads
# timeVsSize(stencilConnection, stencilCursor, destApp, "2019-04-25 10:43:00", "2019-09-24 15:45:00") # 100 threads
# leftDataCumulativeGraph(stencilConnection, stencilCursor, destApp, srcApp)
leftDataCumulativeGraph(stencilConnection, stencilCursor, destApp, srcApp, True)

# g.allTimeVsSizeGraph()

# print len(pd.getMigrationIDsBetweenTimestamps(stencilConnection, stencilCursor, "2019-04-24 11:32:00", "2019-04-24 12:09:00"))
# l = {0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.123, 0.1, 0.1, 0.1, 0.1, 0.2, 0.2, 0.8}
# g.cumulativeGraph(l, "Percentage of Data Left", "Probability")
# print pd.getSizeOfLeftDataInMigratedRows(srcApp, migrationID, stencilCursor)
# print pd.getSizeOfDataWithEntireRowLeft(srcApp, migrationID, stencilCursor)
# print getAllMigrationIDs(stencilConnection, stencilCursor)
# print getMigratedDataSize(destApp, migrationID, stencilCursor)
# print getLengthOfMigration(migrationID, stencilCursor)
db.closeDB(stencilConnection)

