"""
Built and tested in python 3
"""
import pandas
import matplotlib.pyplot as plt
import matplotlib.cm as cm
import numpy as np

from itertools import cycle


df = pandas.read_csv('endpointStats.csv')
df['duration'] /= 1000000 # nanoseconds to ms
df['when'] -= df['when'].max() # relative times

colors = cm.gist_rainbow(np.linspace(0, 1, len(df.ENDPOINT.unique())))
lines = ["-", "--",]
linecycler = cycle(lines)

for endpoint, color in zip(df.ENDPOINT.unique(), colors):
    subset = df.loc[df['ENDPOINT'] == endpoint]
    subset['duration_smooth'] = subset['duration'].rolling(100).mean()

    '''
    plt.scatter(
        subset['when'],
        subset['duration_smooth'],
        label=endpoint,
        color=color,
        s=1,
    )'''
    plt.plot(
        np.unique(subset['when']),
        subset['duration_smooth'],
        color=color,
        label=endpoint,
        linestyle=next(linecycler),
        linewidth=2,
    )

plt.legend(
    loc='best',
    prop={'size': 10},
)
plt.show()
