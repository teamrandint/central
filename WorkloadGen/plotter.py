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
        log=True,
        range=(0,8000)
    )
    plt.show()

def plot_multi_time():
    df1 = pandas.read_csv('./statruns/endpointStats_1.csv')['duration'] / 1000000
    df2 = pandas.read_csv('./statruns/endpointStats_2.csv')['duration'] / 1000000
    df3 = pandas.read_csv('./statruns/endpointStats_5.csv')['duration'] / 1000000
    df4 = pandas.read_csv('./statruns/endpointStats_6.csv')['duration'] / 1000000
    df5 = pandas.read_csv('./statruns/endpointStats_7.csv')['duration'] / 1000000

    
    data = [df1.values,
            df2.values,
            df3.values,
            df4.values,
            df5.values,]

    plt.xlabel('Response time per command (ms)')
    plt.ylabel('Count')
    plt.hist(
        data,
        bins=1000,
        histtype='step',
        stacked=False,
        fill=False,
        log=True,
    )
    plt.show()

def plot_time_stats():

    print("AVERAGE: ", df['duration'].mean())
    for n in range(1, 10):
        print(
            "COMMANDS LESS THAN {}ms:".format(n),
            (df['duration'] < n).mean(),
        )
    for n in range(1, 10):
        print(
            "COMMANDS LESS THAN {}ms:".format(n*10),
            (df['duration'] < n*10).mean()
        )
    for n in range(1, 10):
        print(
            "COMMANDS LESS THAN {}ms:".format(n*100),
            (df['duration'] < n*100).mean()
        )

    for n in range(1, 11):
        print(
            "COMMANDS LESS THAN {}s:".format(n),
            (df['duration'] < n*1000).mean()
        )
    
    print("MAX: ", df['duration'].max())

def plot_tps():
    print(df['when'].min())
    plt.xlabel('Time to last request (s)')
    plt.ylabel('TPS')
    plt.hist(
        df['when'],
        bins=abs(int(df['when'].min() )),
        log=False,
        range=(int(df['when'].min() ), 0)
    )
    plt.show()

def plot_multi_tps():
    df1 = pandas.read_csv('./statruns/endpointStats_1.csv')['when'] 
    df2 = pandas.read_csv('./statruns/endpointStats_2.csv')['when']
    df3 = pandas.read_csv('./statruns/endpointStats_5.csv')['when']
    df4 = pandas.read_csv('./statruns/endpointStats_6.csv')['when']
    df5 = pandas.read_csv('./statruns/endpointStats_7.csv')['when']

    df1 -= df1.min()
    df2 -= df2.min()
    df3 -= df3.min()
    df4 -= df4.min()
    df5 -= df5.min()

    df1 /= 1e9 # nanos to s
    df2 /= 1e9 # nanos to s
    df3 /= 1e9 # nanos to s
    df4 /= 1e9 # nanos to s
    df5 /= 1e9 # nanos to s
    data = [df1.values,
            df2.values,
            df3.values,
            df4.values,
            df5.values,]


    plt.xlabel('Time until last command (s)')
    plt.ylabel('TPS')
    plt.hist(
        data,
        bins=120,
        range=(0, 120),
        histtype='step',
        stacked=False,
        fill=False,
        log=True,
    )
    plt.show()


plot_multi_time()