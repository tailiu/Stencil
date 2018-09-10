import numpy, math
from scipy.linalg import eigh
import matplotlib.pyplot as plt

def getMedian (vectors):
    return numpy.median(vectors)

def getEuclideanDistance (v1, v2):
    return round(numpy.linalg.norm(v1 - v2), 4)

def getGamma(median):
    return round(float(1/median), 4)

def getBandwidthKernelResult(distance):
    return round(math.exp( -1 * ( math.pow(distance, 2) / ( math.pow(0.45, 2) ) ) ), 4)

def getGammaKernelResult(gamma, distance):
    return round(math.exp( -1 * gamma * math.pow(distance, 2) ), 4)

if __name__ == "__main__":

    matrix = numpy.matrix([
        [2.6625, 1.1069],
        [2.2941, 1.5472],
        [0.7372, 0.7908],
        [2.7260, 1.6995],
        [1.1467, 0.8607],
        [2.4963, 1.6729],
        [0.8335, 0.5753],
        [0.7268, 1.1422],
    ])

    distances = numpy.zeros((len(matrix),len(matrix)))
    kernel_mx = numpy.zeros((len(matrix),len(matrix)))
    result_mx = numpy.zeros((len(matrix),1))

    for i in range(0, len(matrix)):
        for j in range(0, len(matrix)):
            distances[i][j] = getEuclideanDistance(matrix[i], matrix[j])

    median = getMedian(distances)
    gamma  = getGamma(median)

    for i in range(0, len(matrix)):
        for j in range(0, len(matrix)):
            kernel_mx[i][j] = getGammaKernelResult(gamma, distances[i][j])
            #kernel_mx[i][j] = getBandwidthKernelResult(distances[i][j])
    
    kernel_mx2 = numpy.matrix([
        [1.0, 0.2423, 0.0081, 0.2295, 0.0226, 0.2330, 0.0091, 0.0084],
        [0.2423, 1.0, 0.0139, 0.3228, 0.0368, 0.5556, 0.0131, 0.0184],
        [0.0081, 0.0139, 1.0, 0.0045, 0.3585, 0.0078, 0.5583, 0.4197],
        [0.2295, 0.3228, 0.0045, 1.0, 0.0121, 0.5650, 0.0044, 0.0059],
        [0.0226, 0.0368, 0.3585, 0.0121, 1.0, 0.0205, 0.3512, 0.2870],
        [0.2330, 0.5556, 0.0078, 0.5650, 0.0205, 1.0, 0.0073, 0.0104],
        [0.0091, 0.0131, 0.5583, 0.0044, 0.3512, 0.0073, 1.0, 0.2406],
        [0.0084, 0.0184, 0.4197, 0.0059, 0.2870, 0.0104, 0.2406, 1.0],
    ])

    eval, evec = eigh(kernel_mx2)

    n_mx = numpy.full([len(matrix), 1], float(1.0/len(matrix)))

    for i in range(0, len(evec)):
        result = numpy.matmul([evec[i]], n_mx)[0][0]
        result_mx[i][0] = round(math.log(math.pow(result, 2) * eval[i]),4)

    print result_mx

    plt.plot(result_mx)
    plt.show()