import graph as g
import proc_data as pd
import numpy as np
import copy

logDir = "../stencil/evaluation/logs/"
leftoverVsMigratedFile = "leftoverVsMigrated"
interruptionTimeFile = "interruptionDuration"
dstAnomalies = "dstAnomaliesVsMigrationSize"
srcAnomalies = "srcAnomaliesVsMigrationSize"
srcSystemDanglingData = "srcSystemDanglingData"
dstSystemDanglingData = "dstSystemDanglingData"
migrationRate = "migrationRate"
dataDownTimeInStencil = "dataDowntimeInStencil"
dataDownTimeInNaive = "dataDowntimeInNaive"
migrationTime = "migrationTime"
migratedDataSize = "migratedDataSize"
dataBags = "dataBags"
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
        for _, line in enumerate(f1):
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

    g.mulLinesDanglingDataCumSum(x, danglingLikesCS, danglingCommentsCS, danglingMessagesCS, danglingTotalCS, danglingStatusesCS, danglingFavCS)

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
        "Stencil",
        "Naive"
    ]
    g.cumulativeGraph(data, labels, "Data Downtime (s)", "Cumulative probability")

def getTime(data):
    time1 = []
    time2 = []
    for i, data1 in enumerate(data):
        if i % 2 == 0:
            time1.append(data1["time"])
        else:
            time2.append(data1["time"])
    return time1, time2

def migrationRateDataModels():
    data1 = readFile3(logDir + migrationTime)
    data2 = readFile3(logDir + migratedDataSize)

    stencilTime, naiveTime = getTime(data1)

    x = []
    for data in data2:
        x.append(data["size"])

    times = [stencilTime, naiveTime]
    sizes = [x, x]
    labels = ["Stencil", "Application"]
    xlabel = 'Migration Time (s)'
    ylabel = 'Migration size (bytes)'

    g.mulPoints(times, sizes, labels, xlabel, ylabel)

def randomWalk():
    data = readFile3(logDir + dataBags)

    # size = 0
    # bags = [size]
    # for i, data1 in enumerate(data):
    #     if i % (len(apps) - 1) == 0 and i != 0:        
    #         size = 0
    #         bags.append(size)
    #     size = size + data1["srcDataBagSize"]
    #     bags.append(size)
    data2 = [0] * (len(apps) - 1)

    for i, data1 in enumerate(data):
        data2[i % (len(apps) - 1)] += data1["srcDataBagSize"]

    # print data2
    for i, data3 in enumerate(data2):
        data2[i] = float(data3) / float(len(data)/4)
    data2.insert(0, 0.0)
    g.dataBag(data2, apps, "Data bag size (bytes)")


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

# leftoverCDF()
# danglingData()
# interruptionTimeCDF()
danglingDataCumSum()
# serviceInterruptionCumSum()
# anomaliesCumSum()
# danglingDataPoints()
# danglingDataSystem()
# migrationRateDifferentNumOfThreads('Consistent/independent migration', migrationRate)
# migrationRateDifferentNumOfThreads('Deletion migration', migrationRate)
# dataDownTime()
# migrationRateDataModels()
# randomWalk()
# danglingDataSystemCombined()