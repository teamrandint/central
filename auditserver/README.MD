# Audit Server

Right now there isn't much in the means of input validiation or error checking, but it is to come.

Multithreading requests to write to the log object as well.

## Endpoints

For each endpoint, pass the information as URI queries.

For example, to post a new userCommand event, the request URI would be as follows:

`localhost:8080/userCommand?server=TRANS&funds=22.33&transactionNum=11&command=BUY`

For all supported params, being surrounded with brackets indicates optional.

### /userCommand

Supported Params are:

- server
- transactionNum
- command
- (username)
- (stockSymbol)
- (filename)
- (funds)

### /quoteServer

Supported Params are:

- server
- transactionNum
- price
- stockSymbol
- username
- quoteServerTime
- cryptokey

### /accountTransaction

Supported Params are:

- server
- transactionNum
- action
- (username)
- (funds)

### /systemEvent

Supported Params are:

- server
- transactionNum
- command
- (username)
- (stockSymbol)
- (filename)
- (funds)

### /errorEvent

Supported Params are:

- server
- transactionNum
- command
- (username)
- (stocksymbol)
- (filename)
- (funds)
- (errormessage)

### /dumpLog

Supported Params are:

- filename
- (username)

## Return Values

Right now the commands just echo the parsed xml. TODO: figure this out