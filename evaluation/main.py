import graph as g
import proc_data as pd
import numpy as np

logDir = "../stencil/evaluation/_logs/"
leftoverVsMigratedFile = "leftoverVsMigrated"
interruptionTimeFile = "interruptionDuration"
dstAnomalies = "dstAnomaliesVsMigrationSize"
srcAnomalies = "srcAnomaliesVsMigrationSize"

def readFile1(filePath):
    data = []
    with open(filePath) as f1:
        for _, line in enumerate(f1):
            if line == "\n":
                continue
            data.append(float(line.rstrip()))
    return data

def readFile2(filePath):
    data = []
    with open(filePath) as f1:
        for _, line in enumerate(f1):
            if line == "\n":
                continue
            e = pd.convertToFloat(line.rstrip().split(","))
            data += e
    return data

def readFile3(filePath):
    data = []
    with open(filePath) as f1:
        for _, line in enumerate(f1):
            if line == "\n":
                continue
            e = pd.convertToJSON(line.rstrip())
            # print type(e)
            if "likes:comments:posts" in e:
                print e["likes:comments:posts"]
    return data

def leftoverCDF():
    data = readFile1(logDir + leftoverVsMigratedFile)
    g.cumulativeGraph(data, "Percentage of data left in Diaspora", "Cumulative probability")

def interruptionTimeCDF():
    data = readFile2(logDir + interruptionTimeFile)
    g.cumulativeGraph(data, "Service interruption time (s)", "Cumulative probability")

def getCumSum():
    arr = [1, 5, 6, 9]
    y = np.cumsum(arr)
    y = np.arange(10, 1011)
    x = np.arange(1, 1001)
    print y
    print x
    readFile3(logDir + srcAnomalies)
    g.line(x, y, "Number of migrated users", "hahah", "good day")

# leftoverCDF()
# interruptionTimeCDF()
getCumSum()
