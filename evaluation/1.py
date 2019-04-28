from matplotlib.ticker import FuncFormatter
import matplotlib.pyplot as plt
import numpy as np

x = np.arange(0, 100, 25)
money = [1.5e5, 2.5e6, 5.5e6, 2.0e7]


def percentage(x, pos):
    'The two args are the value and tick position'
    return '%d%' % (x * 1e-6)


formatter = FuncFormatter(percentage)

fig, ax = plt.subplots()
ax.xaxis.set_major_formatter(formatter)
plt.bar(x, money, width=25, align="edge", edgecolor='white')
plt.xticks(np.append(x, 100))
plt.show()