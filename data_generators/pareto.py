from numpy import *
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns
import sys
# from paretochart import pareto
# %matplotlib inline

def pareto_plot(df, x=None, y=None, title=None, show_pct_y=False, pct_format='{0:.0%}'):
    xlabel = x
    ylabel = y
    tmp = df.sort_values(y, ascending=False)
    x = tmp[x].values
    y = tmp[y].values
    weights = y / y.sum()
    cumsum = weights.cumsum()
    
    fig, ax1 = plt.subplots()
    ax1.bar(x, y)
    ax1.set_xlabel(xlabel)
    ax1.set_ylabel(ylabel)

    ax2 = ax1.twinx()
    ax2.plot(x, cumsum, '-ro', alpha=0.5)
    ax2.set_ylabel('', color='r')
    ax2.tick_params('y', colors='r')
    
    vals = ax2.get_yticks()
    ax2.set_yticklabels(['{:,.2%}'.format(x) for x in vals])

    # hide y-labels on right side
    if not show_pct_y:
        ax2.set_yticks([])
    
    formatted_weights = [pct_format.format(x) for x in cumsum]
    for i, txt in enumerate(formatted_weights):
        ax2.annotate(txt, (x[i], cumsum[i]), fontweight='heavy')    
    
    if title:
        plt.title(title)
    
    plt.tight_layout()
    plt.show()

#Inverse Transform Sampling
def ITS(alpha,Xm,X):
	Y = []
	for i in X:
		if i >= 0:
			y = Xm/((1.0-i)**(1.0/alpha))#Pareto
			Y.append(y)
	return Y


if __name__ == "__main__":

    X_m = 0.1
    alpha = float(sys.argv[1]) # 2
    total = sys.argv[2] # 1000
    X = arange(0,1,1/float(total))
    Result = ITS(alpha,X_m,X)
    Result = str(Result)
    Result = Result.replace(" ", "").strip("[]")
    print Result