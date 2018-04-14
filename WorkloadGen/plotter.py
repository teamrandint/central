"""
Built and tested in python 3
"""
import pandas
import matplotlib.pyplot as plt
import matplotlib.cm as cm
import numpy as np

from itertools import cycle

endpoints = [
    "QUOTE",
    "ADD",
    "DISPLAY_SUMMARY",

    "SELL",
    "CANCEL_SELL",
    "COMMIT_SELL",
    "BUY",
    "CANCEL_BUY",
    "COMMIT_BUY",

    "SET_BUY_TRIGGER",
    "SET_BUY_AMOUNT",
    "CANCEL_SET_BUY",
    "SET_SELL_TRIGGER",
    "SET_SELL_AMOUNT",
    "CANCEL_SET_SELL",
]

postion = {
    "QUOTE":            311,
    "ADD":              311,
    "DISPLAY_SUMMARY":  311,

    "SELL": 312,
    "CANCEL_SELL": 312,
    "COMMIT_SELL": 312,
    "BUY": 312,
    "CANCEL_BUY": 312,
    "COMMIT_BUY": 312,

    "SET_BUY_TRIGGER": 313,
    "SET_BUY_AMOUNT": 313,
    "CANCEL_SET_BUY": 313,
    "SET_SELL_TRIGGER": 313,
    "SET_SELL_AMOUNT": 313,
    "CANCEL_SET_SELL": 313,
}

def getColorStroke(endpoint):
    return {
        "QUOTE":            ('r', '-'),
        "ADD":              ('g', '-'),
        "DISPLAY_SUMMARY":  ('b', '-'),

        "SELL":             ('r', '-'),
        "CANCEL_SELL":      ('r', '--'),
        "COMMIT_SELL":      ('r', '-.'),
        "BUY":              ('b', '-'),
        "CANCEL_BUY":       ('b', '--'),
        "COMMIT_BUY":       ('b', '-.'),

        "SET_BUY_TRIGGER":  ('r', '-'),
        "SET_BUY_AMOUNT":   ('r', '--'),
        "CANCEL_SET_BUY":   ('r', '-.'),
        "SET_SELL_TRIGGER": ('b', '-'),
        "SET_SELL_AMOUNT":  ('b', '--'),
        "CANCEL_SET_SELL":  ('b', '-.'),
    }[endpoint]

df = pandas.read_csv('endpointStats_final_good.csv')
df['duration'] /= 1000000 # nanoseconds to ms
df['when'] -= df['when'].max() # relative times
df['when'] /= 1e9 # nanos to s

def plot_by_command():
    fig, ax = plt.subplots(sharex=True, sharey=True)

    for endpoint in endpoints:
        color, stroke = getColorStroke(endpoint)

        subset = df.loc[df['ENDPOINT'] == endpoint]
        subset['duration_smooth'] = subset['duration']#.rolling(1000).mean()

        #a = subset['when'][subset['when'].duplicated(keep=False)]
        #print(a)
        plt.subplot(postion[endpoint], sharex=ax, sharey=ax)
        plt.plot(
            np.unique(subset['when']),
            subset['duration_smooth'],
            color=color,
            label=endpoint,
            linestyle=stroke,
            linewidth=1,
        ) 

        plt.legend(
            loc='best',
            prop={'size': 10},
        )

    plt.xlabel('Time to last request (s)')
    plt.ylabel('Rolling avg of response time (ms)')
    plt.show()

def plot_by_time():
    plt.xlabel('Response time per command (ms)')
    plt.ylabel('Count')
    plt.hist(
        df['duration'],
        bins=1000,
        log=False,
        range=(0,8000)
    )
    plt.show()

def plot_time_stats():

    print("AVERAGE: ", df['duration'].mean())
    for n in range(1, 10):
        print(
            "% OF COMMANDS LESS THAN {}ms:".format(n),
            (df['duration'] < n).mean()
        )
    for n in range(1, 10):
        print(
            "% OF COMMANDS LESS THAN {}ms:".format(n*10),
            (df['duration'] < n*10).mean()
        )
    for n in range(1, 10):
        print(
            "% OF COMMANDS LESS THAN {}ms:".format(n*100),
            (df['duration'] < n*100).mean()
        )

    for n in range(1, 10):
        print(
            "% OF COMMANDS LESS THAN {}s:".format(n),
            (df['duration'] < n*1000).mean()
        )
    
    print("MAX: ", df['duration'].max())

plot_time_stats()