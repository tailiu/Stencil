import graph as g
import proc_data as pd
import numpy as np

logDir = "../stencil/evaluation/logs/"
leftoverVsMigratedFile = "leftoverVsMigrated"
interruptionTimeFile = "interruptionDuration"
dstAnomalies = "dstAnomaliesVsMigrationSize"
srcAnomalies = "srcAnomaliesVsMigrationSize"
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
    g.cumulativeGraph(data, "Service interruption time (s)", "Cumulative probability")

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

def danglingDataCumSum():
    srcData = readFile3(logDir + srcAnomalies)
    dstData = readFile3(logDir + dstAnomalies)

    danglingLikes, danglingComments, danglingMessages = getDanglingDataInSrc(srcData)
    danglingLikesCS = np.cumsum(danglingLikes)
    danglingCommentsCS = np.cumsum(danglingComments)
    danglingMessagesCS = np.cumsum(danglingMessages)

    danglingStatuses, danglingFav = getDanglingDataInDst(dstData)
    danglingStatusesCS = np.cumsum(danglingStatuses)
    danglingFavCS = np.cumsum(danglingFav)

    x = np.arange(1, cumNum + 1)

    g.mulLinesDanglingData(x, danglingLikesCS, danglingCommentsCS, danglingMessagesCS, danglingStatusesCS, danglingFavCS)

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

# leftoverCDF()
# interruptionTimeCDF()
danglingDataCumSum()
# serviceInterruptionCumSum()
# anomaliesCumSum()
