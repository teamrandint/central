For a database, we chose to use Redis for its scalaiblity, fault tolerance and low overhead. 
In addition, in an iterative course such as this, we value being able to change the schema on the fly.
Since Redis is a key/value store we kept the schema light.
The documentation has key in the header, with a description of the params in the body.
Our redis instance is fully Docker containerized and runs well with the default settings.

### $USERID:Balance

Contains the balance of the user ID. 
Stored as a floating point number for now, however Redis does offer some accuracy guarantees.

#### Functions:
- AddFunds
- GetFunds
- RemoveFunds


### $USERID:Stocks
Redis hash of the stocks a user owns. Stocks are stored as integers.

#### Functions:
- AddStock
- GetStock
- RemoveFunds

### $USERID:SellOrders

Keeps tracks of user's uncomitted sell orders.

#### Functions:
- PushSell
- PopSell

### $USERID:BuyOrders
Keeps tracks of user's uncomitted buy orders.

#### Functions:
- PushBuy
- PopBuy

### $USERID:SellTriggers
Keeps tracks of user's running triggers.

#### Functions:
- AddSellTrigger
- RemoveSellTrigger
- GetSellTrigger

### $USERID:BuyTriggers
Keeps tracks of user's running triggers.

#### Functions:
- AddBuyTrigger
- RemoveBuyTrigger
- GetBuyTrigger

### $USERID:BalanceReserve
Keeps tracks of user's reserve account balance. This holds funds offset for triggers

### $USERID:StocksReserve
Keeps tracks of user's waiting sell triggers balance.

### $USERID:History
Keeps tracks of all user's account transactions.

_Not implemented yet_


## Other functions
### GetUserInfo 
Returns as user's account information

## Running
- Build the docker container
- Expose the proper ports when running (-p exposed:6397)
- Good to go!
