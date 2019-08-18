import graph as g
import proc_data as pd

logDir = "../stencil/evaluation/logs/"
leftoverVsMigratedFile = "leftoverVsMigrated"
interruptionTimeFile = "interruptionDuration"

def readFile1(filePath):
    data = []
    with open(filePath) as f1:
        for _, line in enumerate(f1):
            data.append(float(line.rstrip()))
    return data

def readFile2(filePath):
    data = []
    with open(filePath) as f1:
        for _, line in enumerate(f1):
            e = pd.convertToFloat(line.rstrip().split(","))
            data += e
    return data

def leftoverCDF():
    data = readFile1(logDir + leftoverVsMigratedFile)
    g.cumulativeGraph(data, "Percentage of data left in Diaspora", "Cumulative probability")

def interruptionTimeCDF():
    data = readFile2(logDir + interruptionTimeFile)
    g.cumulativeGraph(data, "Service interruption time (s)", "Cumulative probability")


# leftoverCDF()
interruptionTimeCDF()
