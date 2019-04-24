import matplotlib.pyplot as plt
import numpy as np

def line(x, y, xlabel, ylabel):
    order = np.argsort(x)
    xs = np.array(x)[order]
    ys = np.array(y)[order]

    ax = plt.axes()
    ax.plot(xs, ys)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.show()

def cumulativeGraph(data, xlabel, ylabel):
    data = sorted(data)
    plt.hist(data, cumulative=-1, normed=1, bins=200)
    # plt.xticks([i+0.5 for i in years], years)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.show()