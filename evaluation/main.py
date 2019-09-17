import graph as g
import proc_data as pd
import numpy as np

logDir = "../stencil/evaluation/logs/"
leftoverVsMigratedFile = "leftoverVsMigrated"
interruptionTimeFile = "interruptionDuration"
dstAnomalies = "dstAnomaliesVsMigrationSize"
srcAnomalies = "srcAnomaliesVsMigrationSize"
srcSystemDanglingData = "srcSystemDanglingData"
dstSystemDanglingData = "dstSystemDanglingData"
migrationRate = "migrationRate"
cumNum = 1000

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

def danglingDataSystem():
    srcData = readFile3(logDir + srcSystemDanglingData)
    dstData = readFile3(logDir + dstSystemDanglingData)
    danglingLikes, danglingComments, danglingMessages = getDanglingDataInSrcSystem(srcData)
    danglingStatuses, danglingFav = getDanglingDataInDstSystem(dstData)

    x = np.arange(1, cumNum + 1)

    data = [danglingLikes, danglingComments, danglingMessages, danglingStatuses, danglingFav]
    title = [
        'Dangling likes without posts in Diaspora',
        'Dangling comments without posts in Diaspora',
        'Dangling messages without conversations in Diaspora',
        'Dangling statuses without conversations in Mastodon',
        'Dangling favourites without statuses in Mastodon'
    ]
    xs = 'Num of migration users'
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

def migrationRatePoints():
    data = readFile3(logDir + migrationRate)

    time, size = getMigrationRate(data)
    # g.line(time, size, "Migration time (s)", "Migration size (Bytes)", "Migration rate")
    g.mulPointsDanglingData(time, size, "migration rate")


# leftoverCDF()
# danglingData()
# interruptionTimeCDF()
# danglingDataCumSum()
# serviceInterruptionCumSum()
# anomaliesCumSum()
# danglingDataPoints()
danglingDataSystem()
migrationRatePoints()