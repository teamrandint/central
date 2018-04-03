"""
Built and tested in python 3
"""
import pandas
import matplotlib.pyplot as plt
import matplotlib.cm as cm
import numpy as np


df = pandas.read_csv('endpointStats.csv')


colors = cm.gist_rainbow(np.linspace(0, 1, len(df.ENDPOINT.unique())))


for endpoint, color in zip(df.ENDPOINT.unique(), colors):
    subset = df.loc[df['ENDPOINT'] == endpoint]
    
    subset['duration_smooth'] = subset['duration'].rolling(100).mean()


    plt.scatter(
        subset['when'],
        subset['duration_smooth'],
        label=endpoint,
        color=color,
        s=1,
    )

plt.legend(
    loc='best',
    prop={'size': 8},
)
plt.ylim(ymin=0)
plt.show()
