import graph as g

logDir = "../stencil/evaluation/logs/"
leftoverVsMigratedFile = "leftoverVsMigrated"

def readFile(filePath):
    data = []
    with open(filePath) as f1:
        for _, line in enumerate(f1):
            data.append(float(line.rstrip()))
    return data

def leftoverVsMigrated():
    data = readFile(logDir + leftoverVsMigratedFile)
    g.cumulativeGraph(data, "Percentage of data left in Diaspora", "Cumulative probability")

leftoverVsMigrated()
