import matplotlib.pyplot as plt
import numpy as np
from numpy.polynomial.polynomial import polyfit

# caption font size
plt.rcParams.update({'font.size': 25})

colors = ['g', 'r', 'b', 'c', 'y', 'k', 'm', 'w']
legendFontSize = ['xx-small', 'x-small', 'small', 'medium', 'large', 'x-large', 'xx-large']
legendLoc = ['best', 'upper right', 'upper left', 'upper center', 'center right']
markers = ["o", "v", "s", "*", "+", "<"]
lineStyles = ["solid", "dashed", "dotted",""]

def line(x, y, xlabel, ylabel, title):
    xs, ys = _sortX(x, y)

    ax = plt.axes()
    ax.plot(xs, ys)
    plt.title(title)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.show()

def _sortX(x, y):
    order = np.argsort(x)
    xs = np.array(x)[order]
    ys = np.array(y)[order]
    return xs, ys

def allTimeVsSizeGraph():
    time1 = [9.0, 3.0, 11.0, 7.0, 9.0, 20.0, 5.0, 7.0, 4.0, 17.0, 13.0, 22.0, 5.0, 24.0, 24.0, 3.0, 18.0, 9.0, 9.0, 5.0, 7.0, 15.0, 6.0, 10.0, 9.0, 11.0, 12.0, 21.0, 12.0, 13.0, 5.0, 4.0, 14.0, 54.0, 9.0, 24.0, 10.0, 8.0, 7.0, 6.0, 8.0, 6.0, 8.0, 12.0, 19.0, 20.0, 65.0, 7.0, 11.0, 2.0, 9.0, 6.0, 10.0, 9.0, 7.0, 7.0, 3.0, 6.0, 4.0, 13.0, 7.0, 18.0, 18.0, 12.0, 13.0, 6.0, 11.0, 4.0, 9.0, 9.0, 9.0, 7.0, 10.0, 11.0, 30.0, 4.0, 19.0, 3.0, 8.0, 15.0, 14.0, 13.0, 31.0, 5.0, 23.0, 16.0, 17.0, 13.0, 24.0, 5.0, 5.0, 10.0, 6.0, 13.0, 19.0, 11.0, 6.0, 7.0, 52.0, 14.0, 12.0, 3.0, 6.0, 18.0, 15.0, 11.0, 8.0, 19.0, 3.0, 6.0, 13.0, 4.0, 30.0, 5.0, 6.0, 11.0, 13.0, 8.0, 27.0, 6.0, 7.0, 6.0, 6.0, 12.0, 5.0, 3.0, 12.0, 14.0, 16.0, 8.0, 6.0, 20.0, 5.0, 14.0, 6.0, 15.0, 10.0, 17.0, 13.0, 8.0, 12.0, 13.0, 8.0, 9.0, 12.0, 15.0, 15.0, 8.0, 3.0, 6.0, 16.0, 10.0, 7.0, 33.0, 9.0, 12.0, 13.0, 39.0, 10.0, 6.0, 22.0, 65.0, 8.0, 6.0, 11.0, 9.0, 3.0, 29.0, 11.0, 13.0, 12.0, 10.0, 5.0, 10.0, 6.0, 16.0, 10.0, 7.0, 9.0, 21.0, 16.0, 44.0, 60.0, 10.0, 10.0, 11.0, 10.0, 19.0, 22.0, 5.0, 8.0, 91.0, 10.0, 8.0, 25.0, 8.0, 6.0, 4.0, 8.0, 19.0, 5.0, 12.0, 8.0, 19.0, 15.0, 16.0, 24.0, 18.0, 7.0, 16.0, 28.0, 4.0, 18.0, 7.0, 9.0, 13.0, 16.0, 12.0, 15.0, 12.0, 15.0, 20.0, 15.0, 8.0, 9.0, 17.0, 16.0, 13.0, 9.0, 7.0, 19.0, 8.0, 5.0, 6.0, 7.0, 7.0, 7.0, 8.0, 6.0, 4.0, 14.0, 28.0, 16.0, 9.0, 7.0, 9.0, 22.0, 6.0, 8.0, 8.0, 10.0, 17.0, 8.0, 4.0, 7.0, 11.0, 12.0, 8.0, 4.0, 12.0, 14.0, 11.0, 13.0, 16.0, 4.0, 23.0, 4.0, 16.0, 12.0, 12.0, 5.0, 12.0, 11.0, 6.0, 12.0, 30.0, 7.0, 8.0, 5.0, 14.0, 13.0, 6.0, 10.0, 10.0, 9.0, 8.0, 21.0, 14.0, 16.0, 9.0, 8.0, 6.0, 13.0, 14.0, 14.0, 27.0, 9.0, 20.0, 5.0, 10.0, 6.0, 9.0, 40.0, 19.0, 14.0, 8.0, 5.0, 12.0, 8.0, 34.0, 96.0, 14.0, 4.0, 11.0, 6.0, 6.0, 8.0, 19.0, 10.0, 8.0, 11.0, 13.0, 6.0, 4.0, 6.0, 5.0, 3.0, 17.0, 88.0, 4.0, 8.0, 6.0, 5.0, 49.0, 7.0, 19.0, 20.0, 1.0, 12.0, 6.0, 11.0, 12.0, 15.0, 3.0, 19.0, 5.0, 6.0, 25.0, 11.0, 13.0, 10.0, 8.0, 9.0, 15.0, 13.0, 6.0, 11.0, 21.0, 12.0, 8.0, 5.0, 14.0, 10.0, 12.0, 13.0, 37.0, 6.0, 14.0, 14.0, 7.0, 8.0, 3.0]
    size1 = [15.393, 8.286, 17.234, 13.489, 0, 0, 0, 0, 1.725, 18.085, 0, 19.365, 0, 0, 22.989, 0.633, 0, 0, 0, 0, 0, 20.112, 12.195, 15.516, 0, 0, 16.462, 0, 0, 17.778, 9.503, 2.776, 6.524, 59.646, 0, 22.422, 0, 0, 11.178, 12.643, 0, 0, 0, 14.777, 0, 0, 70.798, 0.642, 0, 0, 11.906, 0, 15.606, 0, 7.597, 0, 0, 0, 1.511, 17.191, 2.497, 18.654, 20.182, 0, 16.021, 0, 13.397, 0, 13.629, 14.557, 0, 0, 11.202, 13.458, 35.738, 0, 14.463, 0, 9.003, 0, 16.067, 0, 32.75, 0, 0, 0, 19.163, 15.969, 0, 3.908, 4.793, 0, 0, 14.357, 0, 0, 0, 0, 73.075, 19.205, 11.059, 0.914, 4.354, 19.664, 19.009, 15.678, 11.25, 19.84, 0.912, 16.666, 18.721, 10.189, 33.96, 4.45, 3.191, 12.594, 19.776, 13.312, 29.156, 0, 0, 0.851, 0, 0, 0, 0, 0, 19.84, 21.604, 0, 0, 0, 0, 0, 2.259, 17.979, 0, 20.571, 19.077, 0, 0, 17.213, 13.774, 0, 0, 14.896, 25.742, 0, 0, 0, 17.724, 14.031, 6.672, 0, 0, 15.871, 17.327, 47.52, 0, 12.097, 21.81, 71.175, 11.77, 0, 0, 8.597, 1.438, 29.394, 0, 20.07, 15.717, 0, 0, 11.127, 0, 18.13, 14.036, 0, 13.174, 0, 0, 39.744, 73.608, 0, 3.885, 0, 0, 0, 23.578, 11.447, 0, 83.062, 3.237, 0, 0, 0.687, 11.651, 3.468, 0, 0, 2.891, 0, 0, 18.195, 19.696, 16.99, 59.4, 14.6, 0, 23.065, 0, 9.87, 21.527, 0, 0, 18.871, 11.458, 15.475, 16.611, 0, 13.862, 0, 17.997, 10.334, 0, 21.217, 18.085, 0, 0, 0, 24.928, 0, 2.507, 0, 0, 0, 0, 16.153, 10.363, 0, 15.547, 0, 17.957, 0, 0, 13.421, 20.802, 0, 0, 3.332, 0, 18.689, 0, 0, 0, 11.395, 13.499, 12.903, 0, 16.096, 12.708, 0, 18.086, 0, 0, 0, 9.663, 17.847, 18.521, 14.146, 0, 19.382, 0, 10.105, 0, 30.117, 1.948, 0, 0, 17.777, 16.049, 0, 13.94, 0, 0, 0, 20.332, 0, 18.835, 6.826, 0, 0, 15.287, 19.586, 18.319, 30.517, 0, 0, 0, 0, 0, 0, 35.066, 19.818, 19.464, 2.31, 2.656, 0, 13.499, 39.669, 110.735, 6.689, 11.179, 6.151, 4.63, 0, 0, 20.445, 0, 11.894, 13.317, 16.647, 0, 0, 0, 0, 8.862, 0, 92.252, 11.339, 0, 0, 0, 46.177, 12.176, 23.374, 21.586, 0, 12.905, 0, 17.137, 0, 16.365, 2.64, 0, 0, 11.081, 24.867, 11.086, 4.924, 0, 0, 0, 0, 13.827, 0, 0, 0, 12.191, 0, 0, 16.002, 13.346, 16.873, 0, 30.139, 1.609, 13.948, 0, 8.377, 8.432, 1.006]
    time5 = [34.0, 23.0, 24.0, 150.0, 32.0, 13.0, 20.0, 20.0, 21.0, 24.0, 15.0, 48.0, 17.0, 38.0, 22.0, 15.0, 25.0, 30.0, 9.0, 137.0, 31.0, 19.0, 21.0, 19.0, 19.0, 29.0, 23.0, 25.0, 48.0, 17.0, 28.0, 90.0, 22.0]
    size5 = [95.907, 41.265, 58.29, 32.492, 44.986, 27.641, 39.193, 48.605, 48.648, 49.749, 36.764, 34.164, 40.405, 96.933, 47.753, 33.725, 42.053, 57.257, 25.831, 422.452, 59.964, 40.818, 42.486, 35.988, 41.465, 47.861, 47.952, 63.748, 116.886, 37.519, 86.474, 272.928, 41.147]
    time10 = [12.0, 23.0, 19.0, 21.0, 20.0, 26.0, 169.0, 23.0, 7.0, 26.0, 17.0, 88.0, 28.0, 15.0, 32.0, 13.0, 27.0, 21.0, 32.0, 20.0, 21.0, 41.0, 16.0, 11.0, 29.0, 19.0, 23.0, 26.0, 25.0, 15.0, 20.0, 31.0, 30.0, 30.0, 15.0, 25.0, 31.0, 18.0, 13.0, 35.0, 11.0, 24.0, 25.0, 14.0, 72.0, 21.0, 18.0]
    size10 = [31.041, 68.524, 44.397, 58.553, 75.691, 81.437, 71.899, 77.584, 27.834, 72.646, 63.287, 71.227, 74.669, 76.551, 81.175, 29.301, 73.077, 53.395, 68.428, 52.869, 70.972, 111.495, 45.592, 26.826, 88.091, 69.695, 82.592, 64.639, 64.197, 49.479, 65.438, 65.642, 83.034, 63.013, 78.999, 76.937, 92.042, 68.906, 48.458, 87.899, 37.698, 83.956, 96.677, 44.362, 232.932, 53.766, 52.214]
    time20 = [15.0, 19.0, 34.0, 26.0, 28.0, 19.0, 30.0, 25.0, 29.0, 35.0, 30.0, 23.0, 49.0, 37.0, 22.0, 18.0, 35.0, 56.0, 25.0, 33.0, 39.0, 31.0, 25.0, 15.0, 33.0, 30.0, 14.0, 41.0, 23.0, 22.0, 33.0, 43.0, 52.0, 32.0, 16.0, 29.0, 39.0]
    size20 = [69.183, 75.598, 108.264, 89.473, 64.494, 68.87, 95.124, 113.605, 103.342, 97.62, 107.476, 91.052, 94.362, 97.126, 61.361, 55.099, 114.561, 155.053, 50.327, 98.267, 124.004, 123.421, 87.25, 33.076, 96.094, 77.739, 70.433, 80.978, 73.795, 46.224, 111.104, 117.368, 125.679, 79.785, 53.69, 99.498, 119.092]
    time50 = [27.0, 52.0, 59.0, 37.0, 51.0, 45.0, 35.0, 58.0, 33.0, 36.0, 50.0, 46.0, 45.0, 43.0, 39.0, 38.0, 33.0, 59.0, 52.0, 49.0, 44.0, 39.0, 46.0, 43.0, 47.0, 65.0, 40.0, 55.0, 43.0, 56.0, 39.0, 23.0, 42.0, 44.0]
    size50 = [30.951, 106.624, 171.421, 63.501, 155.913, 103.731, 104.741, 85.585, 72.551, 113.563, 129.329, 161.271, 127.562, 127.597, 95.924, 87.029, 82.599, 108.92, 150.595, 101.914, 94.587, 124.269, 104.568, 126.174, 116.509, 131.7, 101.499, 133.518, 115.901, 136.263, 74.569, 46.985, 126.102, 96.857]
    time100 = [239.0, 145.0, 365.0, 216.0, 321.0, 272.0, 99.0, 291.0, 60.0, 278.0, 194.0, 328.0, 165.0]
    size100 = [125.168, 76.55, 112.667, 70.996, 133.86, 115.32, 51.697, 91.589, 53.582, 100.272, 65.797, 87.995, 74.439]

    times = [time1, time5, time10, time20, time50, time100]
    sizes = [size1, size5, size10, size20, size50, size100]

    fig, axs = plt.subplots(2, 3, constrained_layout=True)

    titles = {
        0: "(a) 1 Migration Thread",
        1: "(b) 5 Migration Threads",
        2: "(c) 10 Migration Threads",
        3: "(d) 20 Migration Threads",
        4: "(e) 50 Migration Threads",
        5: "(f) 100 Migration Threads"
    }

    for i in range(len(times)):
        if i == 0:
            size1, time1 = _sortX(size1, time1)
            axs[0, 0].plot(size1, time1)
            axs[0, 0].set_title(titles[i])
        elif i == 1:
            size5, time5 = _sortX(size5, time5)
            axs[0, 1].plot(size5, time5)
            axs[0, 1].set_title(titles[i])
        elif i == 2:
            size10, time10 = _sortX(size10, time10)
            axs[0, 2].plot(size10, time10)
            axs[0, 2].set_title(titles[i])
        elif i == 3:
            size20, time20 = _sortX(size20, time20)
            axs[1, 0].plot(size20, time20)
            axs[1, 0].set_title(titles[i])
        elif i == 4:
            size50, time50 = _sortX(size50, time50)
            axs[1, 1].plot(size50, time50)
            axs[1, 1].set_title(titles[i])
        elif i == 5:
            size100, time100 = _sortX(size100, time100)
            axs[1, 2].plot(size100, time100)
            axs[1, 2].set_title(titles[i])

    for ax in axs.flat:
        ax.set(xlabel="Migration Size (KB)", ylabel="Migration Time (s)")

    for time in times:
        print len(time)

    plt.show()

def processCumulativeData(data):
    h, edges = np.histogram(data, density=True, bins=100, )
    h = np.cumsum(h)/np.cumsum(h).max()

    x = edges.repeat(2)[:-1]
    y = np.zeros_like(x)
    y[1:] = h.repeat(2)

    return x, y

def cumulativeGraph(dataArr, labels, xlabel, ylabel):
    # data = sorted(data)
    # plt.hist(data, cumulative=1, normed=1, bins=100, histtype='step')
    # plt.xticks([i+0.5 for i in years], years)

    # plt.grid(True)
    # plt.xlabel(xlabel)
    # plt.ylabel(ylabel)
    # plt.show()
    
    x = []
    y = []

    for data in dataArr:
        x1, y1 = processCumulativeData(data)
        x.append(x1)
        y.append(y1)

    fig, ax = plt.subplots()
    for i, y1 in enumerate(y):
        ax.plot(x[i], y1, colors[i] + lineStyles[i], lw=2, label = labels[i])

    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)
    ax.grid(True)
    
    # This legend setting should be enought for most graphs
    # legend = ax.legend(loc=legendLoc[1], fontsize=legendFontSize[4])

    # I want to change location by coordinates
    # bbox_to_anchor = (x0, y0, width, height) 
    # (x0,y0) are the lower left corner coordinates of the bounding box.
    legend = ax.legend(bbox_to_anchor=(1, 0.95), loc=legendLoc[1], fontsize=legendFontSize[4])

    plt.show()

def barGraph(x, y, xlabel, ylabel, step):
    plt.bar(x, y, width=step, align="edge", edgecolor='white')
    x.append(100)
    plt.xticks(x)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.show()

def mulLinesDanglingData(x, danglingLikes, danglingComments, danglingMessages):
    fig, ax = plt.subplots()
    ax.plot(x, danglingLikes, 'c--', label='Dangling likes without posts in Diaspora')
    ax.plot(x, danglingComments, 'g--', label='Dangling comments without posts in Diaspora')
    ax.plot(x, danglingMessages, 'r--', label='Dangling messages without comments in Diaspora')
    
    ax.grid(True)

    legend = ax.legend(loc=2, fontsize='x-small')

    plt.show()

def mulPointsDanglingData(x, data, labelName):
    fig, ax = plt.subplots()
    ax.plot(x, data, 'r.', label=labelName)
    ax.grid(True)

    legend = ax.legend(loc=1, fontsize='x-small', numpoints=1)
    
    plt.show()

def mulPoints(x, y, labels, xlabel, ylabel):

    fig, ax = plt.subplots()
    
    for i in range(len(x)):
        ax.plot(x[i], y[i], color=colors[i], label=labels[i], markersize=7, marker=markers[i], linestyle=linestyles[-1])

    ax.grid(True)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

    legend = ax.legend(loc=legendLoc[2], fontsize=legendFontSize[3], numpoints=1)
    
    plt.show()

def mulPoints1(x, y, labels, xlabel, ylabel):

    fig, ax = plt.subplots()

    x = np.array(x)
    y = np.array(y)
    b, m = polyfit(x, y, 1)

    # ax.plot(x, y, 'yo', x, m * x + b, '--k', markersize=7)

    ax.plot(x, y, marker='o', color=colors[0], markersize=7, label=labels, linestyle="")
    ax.plot(x, m * x + b, linestyle='--', color=colors[1], linewidth=2)

    ax.grid(True)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

    legend = ax.legend(loc=legendLoc[3], fontsize=legendFontSize[4], numpoints=1)
    
    plt.show()

def mulPoints2(x, y, labels, xlabel, ylabel):

    fig, ax = plt.subplots()

    ax.plot(x, y, marker='o', color=colors[0], markersize=7, label=labels, linestyle="")

    ax.grid(True)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

    legend = ax.legend(loc=legendLoc[3], fontsize=legendFontSize[4], numpoints=1)
    
    plt.show()

def mulPoints3(x, y, labels, xlabels, ylabels):

    figNum = len(y)

    fig, axs = plt.subplots(nrows=1, ncols=figNum)

    for i, ax in enumerate(axs):
        
        x1 = x[i]
        y1 = y[i]

        for j, x11 in enumerate(x1):

            y11 = y1[j]

            x12 = np.array(x11)
            y12 = np.array(y11)
            b, m = polyfit(x12, y12, 1)

            ax.plot(x12, y12, marker=markers[j], color=colors[j], markersize=7, label=labels[j], linestyle=lineStyles[-1])
            ax.plot(x12, m * x12 + b, linestyle=lineStyles[j], color=colors[j], linewidth=2)

        ax.grid(True)
        ax.set_xlabel(xlabels[i])
        ax.set_ylabel(ylabels[i])

        legend = ax.legend(loc=legendLoc[3], fontsize=legendFontSize[3], numpoints=1)
    
    plt.show()

def mulLines(x, y, labels, xlabel, ylabel):

    fig, ax = plt.subplots()

    for i in range(len(y)):
        ax.plot(x, y[i], color=colors[i], label=labels[i], linewidth=3.3, linestyle=linestyles[i])
    
    ax.grid(True)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)

    legend = ax.legend(loc=legendLoc[3], fontsize=legendFontSize[3], numpoints=1)

    plt.show()

def mulLinesDanglingDataCumSum(x, danglingLikesCS, danglingCommentsCS, 
    danglingMessagesCS, danglingTotalCS, danglingStatusesCS, danglingFavCS):
    
    fig, ax = plt.subplots()
    # ax.plot(x, danglingLikesCS, 'c--', label='Dangling likes without posts in Diaspora')
    # ax.plot(x, danglingCommentsCS, 'g--', label='Dangling comments without posts in Diaspora')
    # ax.plot(x, danglingMessagesCS, 'r--', label='Dangling messages without conversations in Diaspora')
    # ax.plot(x, danglingTotalCS, 'k--', label='Dangling likes, comments and messages in Diaspora')
    ax.plot(x, danglingStatusesCS, 'b', label='Dangling statuses without conversations in Mastodon')
    ax.plot(x, danglingFavCS, 'm', label='Dangling favourites without statuses in Mastodon')

    ax.grid(True)

    legend = ax.legend(loc=2, fontsize='xx-small')

    plt.show()

def mulLinesServiceInterruption(x, likesAfterPostsCS, commentsAfterPostsCS, messagesAfterConversationsCS):
    fig, ax = plt.subplots()
    ax.plot(x, likesAfterPostsCS, 'c-', label='Likes migrated after posts')
    ax.plot(x, commentsAfterPostsCS, 'g-', label='Comments migrated after posts')
    ax.plot(x, messagesAfterConversationsCS, 'r-', label='Messages migrated after conversations')

    ax.grid(True)

    legend = ax.legend(loc=2, fontsize='small')

    plt.show()

def mulLinesAnomalies(x, favBeforeStatusesCS, statusesBeforeParentStatusesCS, statusesBeforeConversationsCS):
    fig, ax = plt.subplots()
    ax.plot(x, favBeforeStatusesCS, 'c-', label='Favourites arrived before statuses')
    ax.plot(x, statusesBeforeParentStatusesCS, 'g-', label='Statuses arrived before their corresponding parent statuses')
    ax.plot(x, statusesBeforeConversationsCS, 'r-', label='Statuses arrived before conversations')

    ax.grid(True)

    legend = ax.legend(loc=2, fontsize='small')

    plt.show()

def dataBag(data, apps, ylabel):
    fig, ax = plt.subplots()

    x = np.arange(0, len(apps))

    for i in range(len(data)/len(apps)):
        ax.plot(x, data[i:i+len(apps)], colors[0] + "-.")
    
    ax.set_xticks(x)
    ax.set_xticklabels(apps, fontsize=18)
    ax.grid(True)
    ax.set_ylabel(ylabel)

    plt.show()

def danglingDataSystemCombined(x, y, xlabel, ylabel, labels):
    fig, ax = plt.subplots()
    for i, y1 in enumerate(y):
        ax.plot(x, y1, colors[i] + lineStyles[0], label=labels[i])
    
    ax.grid(True)
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)
    legend = ax.legend(loc=2, fontsize='x-small')

    plt.show()