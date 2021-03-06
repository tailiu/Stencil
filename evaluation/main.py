import graph as g
import proc_data as pd
import numpy as np
import copy

logDir = "../stencil/evaluation/logs/"
evalDir = "../stencil/evaluation/"
leftoverVsMigratedFile = "leftoverVsMigrated"
interruptionTimeFile = "interruptionDuration"
dstAnomalies = "dstAnomaliesVsMigrationSize"
srcAnomalies = "srcAnomaliesVsMigrationSize"
srcSystemDanglingData = "srcSystemDanglingData"
dstSystemDanglingData = "dstSystemDanglingData"
migrationRate = "migrationRate"
dataDownTimeInStencil = "dataDowntimeInStencil"
dataDownTimeInNaive = "dataDowntimeInNaive"
dataDownTimeInPercentageInStencil = "dataDownTimeInPercentageInStencil"
dataDownTimeInPercentageInNaive = "dataDownTimeInPercentageInNaive"
migrationTime = "migrationTime"
migratedDataSize = "migratedDataSize"
migrationScalabilityEdgeFile = "migrationScalabilityEdges"
migrationScalabilityNodeFile = "migrationScalabilityNodes"
counterFile = "counter"
migratedDataSizeBySrcFile = "migratedDataSizeBySrc"
migratedDataSizeByDstFile = "migratedDataSizeByDst"
migratedTimeBySrcFile = "migrationTimeBySrc"
migratedTimeByDstFile = "migrationTimeByDst"
anomaliesFile1 = "danglingData"
anomaliesFile2 = "danglingObjects"
scalabilityFile = "scalability"
scalabilityWithIndepFile = "scalabilityWithIndependent"
dataBagsEnabledFile = "dataBagsEnabled"
dataBagsNotEnabledFile = "dataBagsNotEnabled"

# dataBagsEnabledFile = "dataBagsEnabled1"
# dataBagsNotEnabledFile = "dataBagsNotEnabled1"

cumNum = 1000
apps = ['Diaspora','Mastodon', 'Twitter', 'GNU Social', 'Diaspora']

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
        for i, line in enumerate(f1):
            # if i >= 30:
            #     break
            if line == "\n":
                continue
            obj = pd.convertToJSON(line.rstrip())
            # print type(e)
            data.append(obj)
    return data

def leftoverCDF():
    data = readFile1(logDir + leftoverVsMigratedFile)
    g.cumulativeGraph(data, "Percentage of data left in Diaspora", "Cumulative probability")

def interruptionTimeCDF():
    data = readFile2(logDir + interruptionTimeFile)
    g.cumulativeGraph(data, "Transient dangling data time (s)", "Cumulative probability")

def returnNumOrZero(data, key):
    if key in data:
        return data[key]
    else:
        return 0

def getDanglingDataInSrc(data):
    danglingLikes = []
    danglingComments = []
    danglingMessages = []
    for i, data1 in enumerate(data):
        if i % 2 == 1:
            danglingLikes.append(returnNumOrZero(data1, "likes:posts"))
            danglingComments.append(returnNumOrZero(data1, "comments:posts"))
            danglingMessages.append(returnNumOrZero(data1, "messages:conversations"))
    return danglingLikes, danglingComments, danglingMessages

def getDanglingDataInDst(data):
    danglingStatuses = []
    danglingFav = []
    for i, data1 in enumerate(data):
        if i % 2 == 1:
            danglingStatuses.append(returnNumOrZero(data1, "statuses.conversation_id:conversations.id"))
            danglingFav.append(returnNumOrZero(data1, "favourites.status_id:statuses.id"))
    return danglingStatuses, danglingFav

def danglingData():
    srcData = readFile3(logDir + srcAnomalies)
    danglingLikes, danglingComments, danglingMessages = getDanglingDataInSrc(srcData)

    x = np.arange(1, cumNum + 1)

    g.mulLinesDanglingData(x, danglingLikes, danglingComments, danglingMessages)

def danglingDataPoints():
    srcData = readFile3(logDir + srcAnomalies)
    danglingLikes, danglingComments, danglingMessages = getDanglingDataInSrc(srcData)

    x = np.arange(1, cumNum + 1)

    data = [danglingLikes, danglingComments, danglingMessages]
    label = [
        'Dangling likes without posts in Diaspora',
        'Dangling comments without posts in Diaspora',
        'Dangling messages without comments in Diaspora'
    ]
    for i in range(len(data)):
        g.mulPointsDanglingData(x, data[i], label[i])

def danglingDataCumSum():
    srcData = readFile3(logDir + srcAnomalies)
    dstData = readFile3(logDir + dstAnomalies)

    danglingLikes, danglingComments, danglingMessages = getDanglingDataInSrc(srcData)
    total = [sum(x) for x in zip(danglingLikes, danglingComments, danglingMessages)]
    danglingTotalCS = np.cumsum(total)
    danglingLikesCS = np.cumsum(danglingLikes)
    danglingCommentsCS = np.cumsum(danglingComments)
    danglingMessagesCS = np.cumsum(danglingMessages)

    danglingStatuses, danglingFav = getDanglingDataInDst(dstData)
    danglingStatusesCS = np.cumsum(danglingStatuses)
    danglingFavCS = np.cumsum(danglingFav)

    x = np.arange(1, cumNum + 1)

    g.mulLinesDanglingDataCumSum(x, danglingLikesCS, danglingCommentsCS, danglingMessagesCS, 
        danglingTotalCS, danglingStatusesCS, danglingFavCS)

def getServiceInterruptionData(data):
    likesAfterPosts = []
    commentsAfterPosts = []
    messagesAfterConversations = []
    for i, data1 in enumerate(data):
        if i % 2 == 0:
            likesAfterPosts.append(returnNumOrZero(data1, "likes.target_id:posts.id"))
            commentsAfterPosts.append(returnNumOrZero(data1, "comments.commentable_id:posts.id"))
            messagesAfterConversations.append(returnNumOrZero(data1, "messages.conversation_id:conversations.id"))
    return likesAfterPosts, commentsAfterPosts, messagesAfterConversations

def serviceInterruptionCumSum():
    data = readFile3(logDir + srcAnomalies)
    
    likesAfterPosts, commentsAfterPosts, messagesAfterConversations = getServiceInterruptionData(data)
    likesAfterPostsCS = np.cumsum(likesAfterPosts)
    commentsAfterPostsCS = np.cumsum(commentsAfterPosts)
    messagesAfterConversationsCS = np.cumsum(messagesAfterConversations)

    x = np.arange(1, cumNum + 1)

    g.mulLinesServiceInterruption(x, likesAfterPostsCS, commentsAfterPostsCS, messagesAfterConversationsCS)

def getAnomaliesData(data):
    favBeforeStatuses = []
    statusesBeforeParentStatuses = []
    statusesBeforeConversations = []

    for i, data1 in enumerate(data):
        if i % 2 == 0:
            favBeforeStatuses.append(returnNumOrZero(data1, "favourites.status_id:statuses.id"))
            statusesBeforeParentStatuses.append(returnNumOrZero(data1, "statuses.in_reply_to_id:statuses.id"))
            statusesBeforeConversations.append(returnNumOrZero(data1, "statuses.conversation_id:conversations.id"))
    return favBeforeStatuses, statusesBeforeParentStatuses, statusesBeforeConversations

def anomaliesCumSum():
    data = readFile3(logDir + dstAnomalies)
    
    favBeforeStatuses, statusesBeforeParentStatuses, statusesBeforeConversations = getAnomaliesData(data)
    favBeforeStatusesCS = np.cumsum(favBeforeStatuses)
    statusesBeforeParentStatusesCS = np.cumsum(statusesBeforeParentStatuses)
    statusesBeforeConversationsCS = np.cumsum(statusesBeforeConversations)

    x = np.arange(1, cumNum + 1)

    g.mulLinesAnomalies(x, favBeforeStatusesCS, statusesBeforeParentStatusesCS, statusesBeforeConversationsCS)

def getDanglingDataInSrcSystem(data):
    danglingLikes = []
    danglingComments = []
    danglingMessages = []
    for i, data1 in enumerate(data):
        danglingLikes.append(returnNumOrZero(data1, "likes:posts"))
        danglingComments.append(returnNumOrZero(data1, "comments:posts"))
        danglingMessages.append(returnNumOrZero(data1, "messages:conversations"))
    return danglingLikes, danglingComments, danglingMessages

def getDanglingDataInDstSystem(data):
    danglingStatuses = []
    danglingFav = []
    for i, data1 in enumerate(data):
        danglingStatuses.append(returnNumOrZero(data1, "statuses:conversations"))
        danglingFav.append(returnNumOrZero(data1, "favourites:statuses"))
    return danglingStatuses, danglingFav

def getPercentageInX():
    x = []
    distribution = np.arange(1, cumNum + 1)
    for i in distribution:
        x.append(float(i)/float(cumNum))
    return x

def danglingDataSystem():
    srcData = readFile3(logDir + srcSystemDanglingData)
    dstData = readFile3(logDir + dstSystemDanglingData)
    danglingLikes, danglingComments, danglingMessages = getDanglingDataInSrcSystem(srcData)
    danglingStatuses, danglingFav = getDanglingDataInDstSystem(dstData)

    x = getPercentageInX()

    data = [danglingLikes, danglingComments, danglingMessages, danglingStatuses, danglingFav]
    title = [
        'Dangling likes without posts in Diaspora',
        'Dangling comments without posts in Diaspora',
        'Dangling messages without conversations in Diaspora',
        'Dangling statuses without conversations in Mastodon',
        'Dangling favourites without statuses in Mastodon'
    ]
    xs = 'Percentage of migrated users'
    ys = [
        'Num of dangling likes',
        'Num of dangling comments',
        'Num of dangling messages',
        'Num of dangling statuses',
        'Num of dangling favourites'
    ]
    for i in range(len(data)):
        g.line(x, data[i], xs, ys[i], title[i])
        # g.mulPointsDanglingData(x, data[i], title[i])

def getMigrationRate(data):
    time = []
    size = []
    for i, data1 in enumerate(data):
        time.append(returnNumOrZero(data1, "time"))
        size.append(returnNumOrZero(data1, "size"))
    return time, size

def migrationRateDifferentNumOfThreads(title, fileName):
    data = []
    data.append(readFile3(logDir + fileName + "_1"))
    data.append(readFile3(logDir + fileName + "_10"))
    data.append(readFile3(logDir + fileName + "_50"))
    data.append(readFile3(logDir + fileName + "_100"))

    time = []
    size = []

    for data1 in data:
        time1, size1 = getMigrationRate(data1)
        time.append(time1)
        size.append(size1)
    
    labels = ["1 thread", "10 threads", "50 threads", "100 threads"]
    xlabel = 'Migration Time (s)'
    ylabel = 'Migration size (bytes)'
    title = title
    # g.line(time, size, "Migration time (s)", "Migration size (Bytes)", "Migration rate")
    g.mulPoints(time, size, labels, xlabel, ylabel, title)

def dataDownTime():
    
    data = []
    data.append(readFile2(logDir + dataDownTimeInStencil))
    data.append(readFile2(logDir + dataDownTimeInNaive))
    
    labels = [
        "SA1",
        "Naive system+"
    ]

    xlabel = "Data downtime (s)"
    ylabel = "Cumulative probability"

    g.cumulativeGraph(data, labels, xlabel, ylabel)

def mul100(data):
    res = []

    for data1 in data:
        res.append(data1 * 100.0)

    return res

def dataDownTimeInPercentages(labels):
    
    data = []
    data.append(readFile2(logDir + dataDownTimeInPercentageInStencil))
    data.append(readFile2(logDir + dataDownTimeInPercentageInNaive))
    
    data[0] = mul100(data[0])
    data[1] = mul100(data[1])

    xlabel = "Percentage of time when data is unavailable"
    ylabel = "Cumulative probability"

    # print data

    g.cumulativeGraph(data, labels, xlabel, ylabel)

def convertBytesToMB(data):
    
    for i, data1 in enumerate(data):
        data[i] = float(data1) / 1000000.0
    
    return data

def danglingDataSystemCombined():
    
    srcData = readFile3(logDir + srcSystemDanglingData)
    danglingLikes, danglingComments, danglingMessages = getDanglingDataInSrcSystem(srcData)

    x = getPercentageInX()

    print len(x)
    print len(danglingLikes)
    
    data = [danglingLikes, danglingComments]
    
    xlabel = 'Percentage of migrated users'
    ylabel = 'Dangling data size'
    labels = [
        'likes',
        'comments'
    ]

    g.danglingDataSystemCombined(x, data, xlabel, ylabel, labels)

def getTimeGroups(data, groupNum):

    times = []
    
    for i in range(groupNum):
        times.append([])
    
    for i, data1 in enumerate(data):
        group = i % groupNum
        times[group].append(float(data1["time"]))
    
    return times

def migrationRate(labels):

    times = readFile3(logDir + migrationTime)
    sizes = readFile3(logDir + migratedDataSize)

    groupNum = len(labels)

    times = getTimeGroups(times, groupNum)

    x = []
    for i, data in enumerate(sizes):
        if i % groupNum == 0:
            x.append(data["size"])
    
    x = convertBytesToMB(x)
    sizes = [x] * groupNum

    xlabel = 'Migration size (MB)'
    ylabel = 'Migration time (s)'

    g.mulPoints(sizes, times, labels, xlabel, ylabel)

def getTimeGroups1(data, groupNum):

    times = []
    
    for i in range(groupNum):
        times.append([])
    
    for i, data1 in enumerate(data):
        key = "deletion_time"
        times[0].append(float(data1[key]))
        key = "naive_time"
        times[1].append(float(data1[key]))

    return times

def migrationRate1(labels):

    times = readFile3(logDir + migrationTime)
    sizes = readFile3(logDir + migratedDataSize)

    groupNum = len(labels)

    times = getTimeGroups1(times, groupNum)

    x = []
    for i, data in enumerate(sizes):
        x.append(data["size"])
    
    x = convertBytesToMB(x)
    sizes = [x] * groupNum

    xlabel = 'Migration size (MB)'
    ylabel = 'Migration time (s)'

    g.mulPoints(sizes, times, labels, xlabel, ylabel)

def readDataFromFile(fileName):
    return readFile3(logDir + fileName)
    
def getDataByKey1(data, keyName):
    return float(data[keyName])

def migrationRate2(sizeFiles, timeFiles, labels):

    NUM = 100

    sizeData = []
    timeData = []

    for sizeFile in sizeFiles:
        sizeData.append(readDataFromFile(sizeFile))

    for timeFile in timeFiles:
        timeData.append(readDataFromFile(timeFile))

    groupNum = len(timeData)

    sizes = [[] for _ in range(groupNum)]
    times = [[] for _ in range(groupNum)]

    for group, size in enumerate(sizeData):
        for i, sizeData1 in enumerate(size):
            if i < NUM:
                sizes[group].append(getDataByKey1(sizeData1, "size"))

    for group, time in enumerate(timeData):
        for i, timeData1 in enumerate(time):
            if i < NUM:
                times[group].append(getDataByKey1(timeData1, "time")) 
    
    sizesMB = convertBytesToMB(sizes[0])
    sizes = [sizesMB] * groupNum

    print len(sizes)
    print len(times)
    print len(labels)

    xlabel = 'Migration size (MB)'
    ylabel = 'Migration time (s)'

    g.mulPoints(sizes, times, labels, xlabel, ylabel)

def getTimeFromData(data):

    times = []
    
    for i in range(len(data)):
        times.append([])

    for i, data1 in enumerate(data):
        for data2 in data1:
            times[i].append(float(data2["time"]))
    
    return times

def getTimeFromData1(data):

    times = []
    
    for data1 in data:
        times.append(float(data1["time"]))
    
    return times

def getSizeFromData(data):

    sizes = []
    
    for i in range(len(data)):
        sizes.append([])

    for i, data1 in enumerate(data):
        for data2 in data1:
            sizes[i].append(float(data2["size"]))
    
    return sizes

def migrationRateDatasetsFig(folders, labels):
    
    data1 = []
    data2 = []

    groupNum = len(labels)

    for folder in folders:
        data1.append(readFile3(evalDir + folder + migratedDataSize))
        data2.append(readFile3(evalDir + folder + migrationTime))

    sizes = getSizeFromData(data1)
    times = getTimeFromData(data2)

    for i, size in enumerate(sizes):
        sizes[i] = convertBytesToMB(size)

    xlabel = 'Migration size (MB)'
    ylabel = 'Migration time (s)'

    g.mulPoints(sizes, times, labels, xlabel, ylabel)

def scalabilityEdge(labels):

    data = readFile3(logDir + migrationScalabilityEdgeFile)
    
    edges = []
    times = []

    for data1 in data:
        edges.append(float(data1["edges"]))
        times.append(float(data1["time"]))
 
    xlabel = 'Edges'
    ylabel = 'Migration time (s)'

    g.mulPoints1(edges, times, labels, xlabel, ylabel)

def scalabilityNode(labels):

    data = readFile3(logDir + migrationScalabilityNodeFile)
    
    nodes = []
    times = []

    for data1 in data:
        nodes.append(float(data1["nodes"]))
        times.append(float(data1["time"]))
 
    xlabel = 'Nodes'
    ylabel = 'Migration time (s)'

    g.mulPoints1(nodes, times, labels, xlabel, ylabel)

def scalability(labels):

    data = readFile3(logDir + scalabilityWithIndepFile)
    
    edgesBeforeMigration = []
    nodesBeforeMigration = []
    edgesAfterMigration = []
    nodesAfterMigration = []
    deletionMigrationTimes = []
    independentMigrationTimes = []
    displayTimes = []

    for data1 in data:
        nodesBeforeMigration.append(int(data1["nodes"]))
        nodesAfterMigration.append(int(data1["nodesAfterMigration"]))
        edgesBeforeMigration.append(int(data1["edges"]))
        edgesAfterMigration.append(int(data1["edgesAfterMigration"]))
        displayTimes.append(float(data1["displayTime"]))
        deletionMigrationTimes.append(float(data1["migrationTime"]))
        independentMigrationTimes.append(float(data1["indepMigrationTime"]))
    
    x = [[nodesBeforeMigration, nodesBeforeMigration, nodesAfterMigration], 
        [edgesBeforeMigration, edgesBeforeMigration, edgesAfterMigration]]
    
    y = [[deletionMigrationTimes, independentMigrationTimes, displayTimes],
        [deletionMigrationTimes, independentMigrationTimes, displayTimes]]

    # print independentMigrationTimes

    xlabels = ["Nodes", "Edges"]
    ylabels = ['Time (s)', 'Time (s)']

    g.mulPoints3(x, y, labels, xlabels, ylabels)

def counter(labels):

    data = readFile3(logDir + counterFile)
    
    edges = []
    nodes = []

    for data1 in data:
        edges.append(float(data1["edges"]))
        nodes.append(float(data1["nodes"]))
 
    xlabel = 'Edges'
    ylabel = 'Nodes'

    g.mulPoints2(edges, nodes, labels, xlabel, ylabel)

def calSum(data):

    res = []

    for data1 in data:
        
        res1 = 0.0
        
        for i in data1:
            
            res1 += i
        
        res.append(res1)

    return res

def calSum1(data):

    res1 = 0.0
    for data1 in data:
        res1 += data1

    return res1

def migrationRateDatasetsTab(folders, labels):

    data1 = []

    for folder in folders:
        data1.append(readFile3(evalDir + folder + migrationTime))

    times = getTimeFromData(data1)

    timesSum = calSum(times)

    for i, t in enumerate(timesSum):
        print labels[i] + ":"
        print t

def getSizeFromData(data):

    times = []
    
    for i in range(len(data)):
        times.append([])

    for i, data1 in enumerate(data):
        for data2 in data1:
            times[i].append(float(data2["size"]))
    
    return times

def migrationRateDatasetsTab1(fileNames):

    data1 = []

    for file in fileNames:
        data1.append(readFile3(logDir + file))

    times = getTimeFromData(data1)

    timesSum = calSum(times)

    sizes = getSizeFromData(data1)

    sizesSum = calSum(sizes)

    for i, t in enumerate(timesSum):
        print fileNames[i] + ":"
        print "times:" + str(t) + "," + "sizes:" + str(sizesSum[i])

def migrationRateDatasetsTab2(baseFileName):

    seqNum = 20

    data = []

    for i in range(seqNum):
        data1 = readFile3(logDir + baseFileName + str(i))
        data += data1

    times = getTimeFromData1(data)

    timesSum = calSum1(times)

    print "times:" + str(timesSum)

def compareTwoMigratedSizes(labels):

    data1 = readFile3(logDir + migratedDataSizeBySrcFile)
    data2 = readFile3(logDir + migratedDataSizeByDstFile)
    data3 = readFile3(logDir + migratedTimeBySrcFile)
    data4 = readFile3(logDir + migratedTimeByDstFile)

    data5 = [data1, data2]
    data6 = [data3, data4]

    sizes = getSizeFromData(data5)
    times = getTimeFromData(data6)

    for i, size in enumerate(sizes):
        sizes[i] = convertBytesToMB(size)

    xlabel = 'Migration size (MB)'
    ylabel = 'Migration time (s)'

    print sizes
    
    g.mulPoints(sizes, times, labels, xlabel, ylabel)

def getPercentageInX1(maxNum):
    x = []
    distribution = np.arange(1, maxNum + 1)
    for i in distribution:
        x.append(float(i)/float(maxNum))
    return x

def getPercentageInX2(maxNum):
    x = []
    distribution = np.arange(1, maxNum + 1)
    for i in distribution:
        x.append(float(i)/float(maxNum) * 100.0)
    return x

def getDataByKey(data, key):
    res = []

    for data1 in data:
        res.append(float(data1[key]))
    
    return res

def getPercentagesOfData(data, dataSum):

    percentages = []

    for data1 in data:
        percentages.append(float(data1)/float(dataSum))

    return percentages

def getPercentagesOfData1(data, dataSum):

    percentages = []

    for data1 in data:
        percentages.append(float(data1)/float(dataSum) * 100.0)

    return percentages

def danglingDataSizesCumSum1(labels):

    # srcTotalDataSize = 824719093
    # dstTotalDataSize = 806743608

    srcTotalDataSize = 30840457
    dstTotalDataSize = 16916125

    data = readFile3(logDir + anomaliesFile1)

    dstDanglingData = getDataByKey(data, "dstDanglingData")
    srcDanglingData = getDataByKey(data, "srcDanglingData")

    dstDanglingDataCumSum = np.cumsum(dstDanglingData)
    srcDanglingDataCumSum = np.cumsum(srcDanglingData)

    srcDanglingDataCumSumPercentage = getPercentagesOfData1(srcDanglingDataCumSum, srcTotalDataSize)
    dstDanglingDataCumSumPercentage = getPercentagesOfData1(dstDanglingDataCumSum, dstTotalDataSize)

    x = getPercentageInX2(len(data))
    y = [srcDanglingDataCumSumPercentage, dstDanglingDataCumSumPercentage]

    xlabel = 'Migrated user number as the percentage of the total user number'
    ylabel = 'Dangling data size as the percentage \n of the total data size'
    
    g.mulLines(x, y, labels, xlabel, ylabel)

def danglingObjsCumSum2(labels):

    # Total objs count for the 999 users migration
    # srcTotalObjs = 397001
    # dstTotalObjs = 236766

    # Total objs count for the 1000 users migration
    srcTotalObjs = 393308
    # dstTotalObjs = 208942
    dstTotalObjs = 119564

    data = readFile3(logDir + anomaliesFile2)

    dstDanglingObjs = getDataByKey(data, "dstDanglingObjs")
    srcDanglingObjs = getDataByKey(data, "srcDanglingObjs")

    dstDanglingObjsCumSum = np.cumsum(dstDanglingObjs)
    srcDanglingObjsCumSum = np.cumsum(srcDanglingObjs)

    srcDanglingObjsCumSumPercentage = getPercentagesOfData1(srcDanglingObjsCumSum, srcTotalObjs)
    dstDanglingObjsCumSumPercentage = getPercentagesOfData1(dstDanglingObjsCumSum, dstTotalObjs)

    x = getPercentageInX2(len(data))
    y = [srcDanglingObjsCumSumPercentage, dstDanglingObjsCumSumPercentage]

    xlabel = 'Percentage of users migrated'
    ylabel = 'Percentage of dangling objects'
    # ylabel = 'Percentage of dangling objects to total objects'
    
    g.mulLines(x, y, labels, xlabel, ylabel)

def getPercentagesOfData(numerators, denominators):

    rates = []
    
    for i in range(len(numerators)):
        rates.append(numerators[i] / denominators[i] * 100)
    
    return rates

def randomWalk1(apps, labels):
    
    data1 = readFile3(logDir + dataBagsEnabledFile)
    data2 = readFile3(logDir + dataBagsNotEnabledFile)

    totalObjs1 = getDataByKey(data1, "totalObjs")
    danglingObjs1 = getDataByKey(data1, "danglingObjs")

    totalObjs2 = getDataByKey(data2, "totalObjs")
    danglingObjs2 = getDataByKey(data2, "danglingObjs")

    percentages1 = getPercentagesOfData(danglingObjs1, totalObjs1)
    percentages2 = getPercentagesOfData(danglingObjs2, totalObjs2)

    ylabel = 'Percentage of dangling objects to total objects'
    
    percentages = [percentages1, percentages2]

    g.dataBag1(percentages, labels, apps, ylabel)

def groupData(data, group):
    
    res = [0] * group

    for i, data1 in enumerate(data):
        res[i % group] += data1  

    return res

def randomWalk(apps):

    group = len(apps) - 1

    data1 = readFile3(logDir + dataBagsEnabledFile)
    data2 = readFile3(logDir + dataBagsNotEnabledFile)

    print len(data1)
    print len(data2)

    appObjs1 = getDataByKey(data1, "appObjs")
    danglingObjs1 = getDataByKey(data1, "danglingObjs")

    appObjs2 = getDataByKey(data2, "appObjs")
    danglingObjs2 = getDataByKey(data2, "danglingObjs")

    appObjsGrouped1 = groupData(appObjs1, group)
    appObjsGrouped2 = groupData(appObjs2, group)

    danglingObjsGrouped1 = groupData(danglingObjs1, group)
    danglingObjsGrouped2 = groupData(danglingObjs2, group)

    totalObjsGrouped1 = [appObjsGrouped1[i] + danglingObjsGrouped1[i] for i in range(len(appObjsGrouped1))]
    totalObjsGrouped2 = [appObjsGrouped2[i] + danglingObjsGrouped2[i] for i in range(len(appObjsGrouped2))]

    percentages1 = getPercentagesOfData(danglingObjsGrouped1, totalObjsGrouped1)
    percentages2 = getPercentagesOfData(danglingObjsGrouped2, totalObjsGrouped2)

    percentages1.insert(0, 0.0)
    percentages2.insert(0, 0.0)
    
    percentages = [percentages1, percentages2]

    return percentages

def randomWalkBar(apps, ylabel, labels):

    percentages = randomWalk(apps)

    g.dataBagBar(percentages, labels, apps, ylabel)

def randomWalkLine(apps, ylabel, labels):

    percentages = randomWalk(apps)

    print percentages

    g.dataBagLine(apps, percentages, labels, ylabel)

def addData(data, start, step):
    total = 0
    for i in range(start, len(data), step):
        total += data[i]
    return total

def randomWalk1(apps):

    group1 = len(apps)
    group2 = len(apps) * 2

    data1 = readFile3(logDir + dataBagsEnabledFile)
    data2 = readFile3(logDir + dataBagsNotEnabledFile)

    print len(data1)
    print len(data2)

    appObjs1 = getDataByKey(data1, "appObjs")
    appObjs2 = getDataByKey(data2, "appObjs")

    appObjsGrouped1 = []
    appObjsGrouped2 = []
    for i in range(group2):
        appObjsGrouped1.append(addData(appObjs1, i, group2))
        appObjsGrouped2.append(addData(appObjs2, i, group2))
    
    percentages1 = []
    percentages2 = []
    for i in range(group1):
        percentages1.append(float(appObjsGrouped1[i + group1])/float(appObjsGrouped1[i]) * 100.0)
        percentages2.append(float(appObjsGrouped2[i + group1])/float(appObjsGrouped2[i]) * 100.0)

    percentages = [percentages1, percentages2]
    return percentages

def randomWalkLine1(apps, ylabel, labels):

    percentages = randomWalk1(apps)

    print percentages

    g.dataBagLine1(apps, percentages, labels, ylabel)

def dataBags(apps, ylabel, labels):

    data1 = []
    data2 = []
    for app in apps:
        if app == "GNU Social":
            app = "gnusocial"
        data1.append(readFile3(logDir + dataBagsEnabledFile + "_" + app.lower()))
        data2.append(readFile3(logDir + dataBagsNotEnabledFile + "_" + app.lower()))

    appObjs1 = []
    appObjs2 = []
    for data in data1:
        appObjs1.append(getDataByKey(data, "appObjs"))
    for data in data2:
        appObjs2.append(getDataByKey(data, "appObjs"))
    
    appObjsEnabledBefore = []
    appObjsEnabledAfter = []
    for data in appObjs1:
        appObjsEnabledBefore.append(addData(data, 0, 5))
    for data in appObjs1:
        appObjsEnabledAfter.append(addData(data, 4, 5))
    
    appObjsUnenabledBefore = []
    appObjsUnenabledAfter = []
    for data in appObjs2:
        appObjsUnenabledBefore.append(addData(data, 0, 5))
    for data in appObjs2:
        appObjsUnenabledAfter.append(addData(data, 4, 5))

    enabledPercentages = []
    unenabledPercentages = []
    for i in range(len(apps)):
        enabledPercentages.append(float(appObjsEnabledAfter[i])/float(appObjsEnabledBefore[i]) * 100.0)
        unenabledPercentages.append(float(appObjsUnenabledAfter[i])/float(appObjsUnenabledBefore[i]) * 100.0)
    
    print enabledPercentages
    print unenabledPercentages

    percentages = [enabledPercentages, unenabledPercentages]
    g.dataBagBar(percentages, labels, apps, ylabel)

# leftoverCDF()
# danglingData()
# interruptionTimeCDF()
# danglingDataCumSum()
# serviceInterruptionCumSum()
# anomaliesCumSum()
# danglingDataPoints()
# danglingDataSystem()
# migrationRateDifferentNumOfThreads('Consistent/independent migration', migrationRate)
# migrationRateDifferentNumOfThreads('Deletion migration', migrationRate)
# dataDownTime()
# migrationRate(["App with DAG and display", 
    # "App without DAG but with display", 
    # "App with DAG but without display", 
    # "App without DAG or display"])
# randomWalk()
# danglingDataSystemCombined()


# migrationRate1(["SA1", "Naive system"])
# migrationRateDatasetsFig(["logs_1M/", "logs_100K/", "logs_10K/"], ["1M", "100K", "10K"])
# dataDownTime()
# danglingDataSizesCumSum1(["Diaspora (source)", "Mastodon (destination)"])
# scalabilityEdge("SA1")
# scalabilityNode("SA1")
# counter("")
# migrationRateDatasetsTab(["logs_1M/", 
#     "logs_100K/", 
#     "logs_10K/", 
#     "logs_1K/"], 
#     ["1M", "100K", "10K", "1K"])
# migrationRateDatasetsTab1(["diaspora_1K_dataset", "diaspora_10K_dataset", 
#     "diaspora_100K_dataset", "diaspora_1M_dataset"])
# migrationRateDatasetsTab2("diaspora_10K_dataset_sa2_")
# compareTwoMigratedSizes(["Source", "Destination"])
# randomWalk1(["diaspora", "mastodon", "diaspora"],
#             ["Stencil with data bags", "Stencil without data bags"])
# randomWalkBar(["Diaspora", "Mastodon", "Gnusocial", "Twitter", "Diaspora"],
#             'Percentage of dangling objects to total objects',
#             ["Stencil with data bags", "Stencil without data bags"])
# randomWalkLine(["Diaspora", "Mastodon", "Gnusocial", "Twitter", "Diaspora"], 
#             'Percentage of dangling objects to total objects',
#             ["Stencil with \n data bags", "Stencil without \n data bags"])
# randomWalkLine1(["Diaspora", "Mastodon", "Gnusocial", "Twitter"],
#                 'Percentage of application objects',
#                 ["Stencil with data bags", "Stencil without data bags"])


# danglingObjsCumSum2(["Diaspora (source)", "Mastodon (destination)"])
# migrationRate2(["SA1Size", "SA1WDSize", "SA1IndepSize", "naiveSize"], 
#     ["SA1Time", "SA1WDTime", "SA1IndepTime", "naiveTime"], 
#     ["Deletion", "Deletion w/o validation", "Independent", "Naive"])
# scalability(["Deletion", "Independent", "Validation"])
# dataDownTimeInPercentages(["Stencil", "Naive+"])
dataBags(["Diaspora", "Mastodon", "GNU Social", "Twitter"],
        'Percentage of data returned',
        ["With data bags", "Without data bags"])