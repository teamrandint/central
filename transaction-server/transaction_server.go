package main

import (
	"fmt"
	"os"
	"seng468/transaction-server/database"
	"seng468/transaction-server/logger"
	"seng468/transaction-server/quote"
	"seng468/transaction-server/socketserver"
	"seng468/transaction-server/trigger"
	"strconv"

	"errors"

	"github.com/shopspring/decimal"
)

// TransactionServer holds the main components of the module itself
type TransactionServer struct {
	Name          string
	Addr          string
	Server        socketserver.SocketServer
	Logger        logger.Logger
	UserDatabase  database.RedisDatabase
	TriggerClient triggerclient.TriggerClient
}

func main() {
	serverAddr := ":" + os.Getenv("transport")
	databaseAddr := "tcp"
	databasePort := os.Getenv("dbaddr") + ":" + os.Getenv("dbport")
	auditAddr := "http://" + os.Getenv("auditaddr") + ":" + os.Getenv("auditport")
	triggerURL := "http://" + os.Getenv("triggeraddr") + ":" + os.Getenv("triggerport")

	server := socketserver.NewSocketServer(serverAddr)
	database := database.RedisDatabase{Addr: databaseAddr, Port: databasePort, DbRequests: make(chan *database.Query, 1000),
		BatchSize: 20, PollRate: 20, BatchResults: make(chan database.Response, 1000), DbPool: database.NewPool(databaseAddr, databasePort)}
	logger := logger.AuditLogger{Addr: auditAddr}
	triggerclient := triggerclient.TriggerClient{TriggerURL: triggerURL}

	ts := &TransactionServer{
		Name:          "transactionserve",
		Addr:          serverAddr,
		Server:        server,
		Logger:        logger,
		UserDatabase:  database,
		TriggerClient: triggerclient,
	}

	server.Route("ADD,<user>,<amount>", ts.Add)
	server.Route("QUOTE,<user>,<stock>", ts.Quote)
	server.Route("BUY,<user>,<stock>,<amount>", ts.Buy)
	server.Route("COMMIT_BUY,<user>", ts.CommitBuy)
	server.Route("CANCEL_BUY,<user>", ts.CancelBuy)
	server.Route("SELL,<user>,<stock>,<amount>", ts.Sell)
	server.Route("COMMIT_SELL,<user>", ts.CommitBuy)
	server.Route("CANCEL_SELL,<user>", ts.CancelBuy)
	server.Route("SET_BUY_AMOUNT,<user>,<stock>,<amount>", ts.SetBuyAmount)
	server.Route("CANCEL_SET_BUY,<user>,<stock>", ts.CancelSetBuy)
	server.Route("SET_BUY_TRIGGER,<user>,<stock>,<amount>", ts.SetBuyTrigger)
	server.Route("SET_SELL_AMOUNT,<user>,<stock>,<amount>", ts.SetSellAmount)
	server.Route("SET_SELL_TRIGGER,<user>,<stock>,<amount>", ts.SetSellTrigger)
	server.Route("TRIGGER_SUCCESS,<user>,<stock>,<price>,<amount>,<action>", ts.TriggerSuccess)
	server.Route("CANCEL_SET_SELL,<user>,<stock>", ts.CancelSetSell)
	server.Route("DUMPLOG,<user>,<filename>", ts.DumpLogUser)
	server.Route("DISPLAY_SUMMARY,<user>", ts.DisplaySummary)
	go ts.UserDatabase.DbRequestWorker()
	server.Run()
}

// Add the given amount of money to the user's account
// Params: user, amount
// PostCondition: the user's account is increased by the amount of money specified
func (ts TransactionServer) Add(transNum int, params ...string) string {
	user := params[0]
	amount, err := decimal.NewFromString(params[1])
	if err != nil {
		ts.reportError(transNum, "ADD", user, "Could not parse add amount to decimal",
			nil, nil, nil)
		return "-1"
	}

	err = ts.UserDatabase.AddFunds(user, amount)
	if err != nil {
		ts.reportError(transNum, "ADD", user, "Failed to add amount to the database for user: "+err.Error(),
			nil, nil, amount.String())
		return "-1"
	}
	go ts.Logger.AccountTransaction(ts.Name, transNum, "ADD", user, amount)
	return "1"
}

// Quote gets the current quote for the stock for the specified user
// Params: user, stock
// PostCondition: the current price of the specified stock is displayed to the user
func (ts TransactionServer) Quote(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]
	dec, err := quoteclient.Query(user, stock, transNum)
	if err != nil {
		ts.reportError(transNum, "QUOTE", user, err.Error(),
			stock, nil, nil)
		return "-1"
	}
	return dec.StringFixed(2)
}

// Buy the dollar amount of the stock for the specified user at the current price.
// Params: user, stock, amount
// PreCondition: The user's account must be greater or equal to the amount of the purchase.
// PostCondition: The user is asked to confirm or cancel the transaction
func (ts TransactionServer) Buy(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]
	amount, err := decimal.NewFromString(params[2])
	if err != nil {
		ts.reportError(transNum, "BUY", user, "Could not parse buy amount to decimal", stock, nil, nil)
		return "-1"
	}

	curr, err := ts.UserDatabase.GetFunds(user)
	if err != nil {
		ts.reportError(transNum, "BUY", user, fmt.Sprintf("Error connecting to the database to get funds: %s", err.Error()),
			stock, nil, amount.String())
		return "-1"
	}

	if curr.LessThan(amount) {
		ts.reportError(transNum, "BUY", user, "Not enough funds to issue buy order", stock, nil, amount.String())
		return "-1"
	}

	cost, shares, err := ts.getMaxPurchase(user, stock, amount, nil, transNum)
	if err != nil {
		ts.reportError(transNum, "BUY", user, fmt.Sprintf("Error connecting to the quote server: %s", err.Error()),
			stock, nil, amount.String())
		return "-1"
	}

	err = ts.UserDatabase.RemoveFunds(user, cost)
	if err != nil {
		ts.reportError(transNum, "BUY", user, fmt.Sprintf("Error removing funds: %s", err.Error()),
			stock, nil, amount.String())
		return "-1"
	}
	err = ts.UserDatabase.PushBuy(user, stock, cost, shares)
	if err != nil {
		ts.reportError(transNum, "BUY", user, fmt.Sprintf("Error pushing buy command: %s", err.Error()),
			stock, nil, amount.String())
		return "-1"
	}

	go ts.Logger.AccountTransaction(ts.Name, transNum, "remove", user, amount)
	return "1"
}

// CommitBuy commits the most recently executed BUY command
// Params: user
// Pre-Conditions: The user must have executed a BUY command within the previous 60 seconds
// Post-Conditions:
// 		(a) the user's cash account is decreased by the amount user to purchase the stock
// 		(b) the user's account for the given stock is increased by the purchase amount
func (ts TransactionServer) CommitBuy(transNum int, params ...string) string {
	user := params[0]
	go ts.Logger.SystemEvent(ts.Name, transNum, "COMMIT_BUY", user, nil, nil, nil)
	stock, _, shares, err := ts.UserDatabase.PopBuy(user)
	if err != nil {
		ts.reportError(transNum, "COMMIT_BUY", user, "Error popping command in commit buy: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}

	err = ts.UserDatabase.AddStock(user, stock, shares)
	if err != nil {
		ts.reportError(transNum, "COMMIT_BUY", user, "Error connecting to database to add stock: "+err.Error(),
			stock, nil, string(shares))
		return "-1"
	}
	return "1"
}

// CancelBuy cancels the most recently executed BUY Command
// Param: user
// Pre-Condition: The user must have executed a BUY command within the previous 60 seconds
// Post-Condition: The last BUY command is canceled and any allocated system resources are reset and released.
func (ts TransactionServer) CancelBuy(transNum int, params ...string) string {
	user := params[0]
	stock, cost, _, err := ts.UserDatabase.PopBuy(user)
	if err != nil {
		ts.reportError(transNum, "CANCEL_BUY", user, "Error popping command in cancel buy: "+err.Error(),
			nil, nil, nil)
		return "-1"
	}

	if stock == "" {
		ts.reportError(transNum, "CANCEL_BUY", user, "No pending buy orders to pop", nil, nil, nil)
		return "-1"
	}

	err = ts.UserDatabase.AddFunds(user, cost)
	if err != nil {
		ts.reportError(transNum, "CANCEL_BUY", user, "Error connecting to database to add funds: "+err.Error(),
			stock, nil, cost.String())
		return "-1"
	}
	return "1"
}

// Sell the specified dollar mount of the stock currently held by the specified
// user at the current price.
// Param: user, stock, amount
// Pre-condition: The user's account for the given stock must be greater than
// 		or equal to the amount being sold.
// Post-condition: The user is asked to confirm or cancel the given transaction
func (ts TransactionServer) Sell(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]
	amount, err := decimal.NewFromString(params[2])
	if err != nil {
		ts.reportError(transNum, "SELL", user, "Could not parse sell amount to decimal", stock, nil, nil)
		return "-1"
	}
	cost, shares, err := ts.getMaxPurchase(user, stock, amount, nil, transNum)
	if err != nil {
		ts.reportError(transNum, "SELL", user, "Could not connect to the quote server: "+err.Error(),
			stock, nil, amount.String())
		return "-1"
	}

	curr, err := ts.UserDatabase.GetStock(user, stock)
	if curr < shares {
		ts.reportError(transNum, "SELL", user, "Cannot sell more stock than you own", stock,
			nil, amount.String())
		return "-1"
	}

	err = ts.UserDatabase.RemoveStock(user, stock, shares)
	if err != nil {
		ts.reportError(transNum, "SELL", user, "Error removing stock from database: "+err.Error(), stock, nil,
			string(shares))
		return "-1"
	}

	err = ts.UserDatabase.PushSell(user, stock, cost, shares)
	if err != nil {
		ts.reportError(transNum, "SELL", user, "Error pushing sell command to database: "+err.Error(),
			stock, nil, amount.String())
		return "-1"
	}
	return "-1"
}

// CommitSell commits the most recently executed SELL command
// Params: user
// Pre-Conditions: The user must have executed a SELL command within the previous 60 seconds
// Post-Conditions:
// 		(a) the user's account for the given stock is decremented by the sale amount
// 		(b) the user's cash account is increased by the sell amount
func (ts TransactionServer) CommitSell(transNum int, params ...string) string {
	user := params[0]
	go ts.Logger.SystemEvent(ts.Name, transNum, "COMMIT_SELL", user, nil, nil, nil)

	stock, cost, _, err := ts.UserDatabase.PopSell(user)
	if err != nil {
		ts.reportError(transNum, "COMMIT_SELL", user, "Error connecting to database to pop command: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}

	err = ts.UserDatabase.AddFunds(user, cost)
	if err != nil {
		ts.reportError(transNum, "COMMIT_SELL", user, "Error connecting to database to add funds: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}
	return "1"

}

// CancelSell cancels the most recently executed SELL Command
// Params: user
// Pre-conditions: The user must have executed a SELL command within the previous 60 seconds
// Post-conditions: The last SELL command is canceled and any allocated system resources are reset and released.
func (ts TransactionServer) CancelSell(transNum int, params ...string) string {
	user := params[0]
	stock, _, shares, err := ts.UserDatabase.PopSell(user)
	if err != nil {
		ts.reportError(transNum, "CANCEL_SELL", user, "Error connecting to database to pop command: "+err.Error(),
			nil, nil, nil)
		return "-1"
	}

	err = ts.UserDatabase.AddStock(user, stock, shares)
	if err != nil {
		ts.reportError(transNum, "CANCEL_SELL", user, "Error connecting to database to add stock: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}
	return "1"
}

// SetBuyAmount sets a defined amount of the given stock to buy when the
// current stock price is less than or equal to the BUY_TRIGGER
// Params: user, stock, amount
// Pre-condition: The user's cash account must be greater than or equal to the
//		BUY amount at the time the transaction occurs
// Post-condition:
// 		(a) a reserve account is created for the BUY transaction to hold the
//			specified amount in reserve for when the transaction is triggered
// 		(b) the user's cash account is decremented by the specified amount
// 		(c) when the trigger point is reached the user's stock account is
//			updated to reflect the BUY transaction.
func (ts TransactionServer) SetBuyAmount(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]
	amount, err := decimal.NewFromString(params[2])
	if err != nil {
		ts.reportError(transNum, "SET_BUY_AMOUNT", user, "Could not parse set buy amount to decimal",
			stock, nil, nil)
		return "-1"
	}

	curr, err := ts.UserDatabase.GetFunds(user)
	if err != nil {
		ts.reportError(transNum, "SET_BUY_AMOUNT", user, "Could not get funds from database: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}

	if curr.LessThan(amount) {
		ts.reportError(transNum, "SET_BUY_AMOUNT", user, "Not enough funds to execute command", stock,
			nil, amount.String())
		return "-1"
	}

	err = ts.UserDatabase.RemoveFunds(user, amount)
	if err != nil {
		ts.reportError(transNum, "SET_BUY_AMOUNT", user, "Error removing funds from database: "+err.Error(),
			stock, nil, amount.String())
		return "-1"
	}

	err = ts.UserDatabase.AddReserveFunds(user, amount)
	if err != nil {
		// TODO: add funds back into database
		ts.reportError(transNum, "SET_BUY_AMOUNT", user, "Error adding funds to reserve:  "+err.Error(),
			stock, nil, amount.String())
		return "-1"
	}

	err = ts.TriggerClient.SetNewBuyTrigger(transNum, user, stock, amount)
	if err != nil {
		ts.reportError(transNum, "SET_BUY_AMOUNT", user, "Error setting a new buy trigger: "+err.Error(),
			stock, nil, amount.String())
		return "-1"
	}
	// TODO: add trigger to database
	return "1"
}

// CancelSetBuy cancels a SET_BUY command issued for the given stock
// Params: user, stock
// The must have been a SET_BUY Command issued for the given stock by the user
// Post-condition:
// 		(a) All accounts are reset to the values they would have had had the
//			SET_BUY Command not been issued
// 		(b) the BUY_TRIGGER for the given user and stock is also canceled.
func (ts TransactionServer) CancelSetBuy(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]

	cancelled, err := ts.TriggerClient.CancelBuyTrigger(transNum, user, stock)
	if err != nil {
		ts.reportError(transNum, "CANCEL_SET_BUY", user, "Error cancelling a trigger: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}

	err = ts.UserDatabase.RemoveReserveFunds(user, cancelled.GetAmount())
	if err != nil {
		ts.reportError(transNum, "CANCEL_SET_BUY", user, "Error removing funds from reserve: "+err.Error(),
			stock, nil, cancelled.GetCost().String())
		return "-1"
	}

	err = ts.UserDatabase.AddFunds(user, cancelled.GetAmount())
	if err != nil {
		ts.reportError(transNum, "CANCEL_SET_BUY", user, "Error adding funds: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}

	return "1"
}

// SetBuyTrigger sets the trigger point base on the current stock price when
// any SET_BUY will execute.
// Params: user, stock, amount
// Pre-conditions: The user must have specified a SET_BUY_AMOUNT prior to
//		 setting a SET_BUY_TRIGGER
// Post-conditions: The set of the user's buy triggers is updated to
//		include the specified trigger
func (ts TransactionServer) SetBuyTrigger(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]
	triggerAmount, err := decimal.NewFromString(params[2])
	if err != nil {
		go ts.Logger.SystemError(ts.Name, transNum, "SET_BUY_TRIGGER", user, stock, nil, nil,
			"Could not parse set buy trigger amount to decimal")
		return "-1"
	}

	_, err = ts.TriggerClient.StartNewBuyTrigger(transNum, user, stock, triggerAmount)
	if err != nil {
		ts.reportError(transNum, "SET_BUY_TRIGGER", user, "No existing buy trigger for this user and stock",
			stock, nil, triggerAmount.String())
		return "-1"
	}
	return "1"
}

// SetSellAmount sets a defined amount of the specified stock to sell when
// the current stock price is equal or greater than the sell trigger point
// Params: user, stock, amount
// Pre-conditions: The user must have the specified amount of stock in their
//		account for that stock.
// Post-conditions: A trigger is initialized for this username/stock symbol
//		combination, but is not complete until SET_SELL_TRIGGER is executed.
func (ts TransactionServer) SetSellAmount(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]
	amount, err := strconv.ParseInt(params[2], 10, 64)
	if err != nil {
		ts.reportError(transNum, "SET_SELL_AMOUNT", user, "Could not parse set sell amount to decimal",
			stock, nil, nil)
		return "-1"
	}

	curr, err := ts.UserDatabase.GetStock(user, stock)
	if err != nil {
		ts.reportError(transNum, "SET_SELL_AMOUNT", user, "Could not get stock from database: "+err.Error(),
			stock, nil, string(amount))
		return "-1"
	}

	if amount > curr {
		ts.reportError(transNum, "SET_SELL_AMOUNT", user, "Cannot set sell trigger for more stock than you own",
			stock, nil, string(amount))
		return "-1"
	}

	err = ts.TriggerClient.SetNewSellTrigger(transNum, user, stock, amount)
	if err != nil {
		ts.reportError(transNum, "SET_SELL_AMOUNT", user, "Failed to make new sell trigger: "+err.Error(),
			stock, nil, string(amount))
		return "-1"
	}
	return "1"
}

// SetSellTrigger sets the stock price trigger point for executing any
// SET_SELL triggers associated with the given stock and user
// Params: user, stock, amount
// Pre-Conditions: The user must have specified a SET_SELL_AMOUNT prior to
//		setting a SET_SELL_TRIGGER
// Post-Conditions:
// 		(a) a reserve account is created for the specified amount of the
//			given stock
// 		(b) the user account for the given stock is reduced by the max number
//			of stocks that could be purchased and
// 		(c) the set of the user's sell triggers is updated to include the
//			specified trigger.
func (ts TransactionServer) SetSellTrigger(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]
	price, err := decimal.NewFromString(params[2])
	if err != nil {
		ts.reportError(transNum, "SET_SELL_TRIGGER", user, "Could not parse set sell trigger price to decimal",
			stock, nil, nil)
		return "-1"
	}

	trig, err := ts.TriggerClient.StartNewSellTrigger(transNum, user, stock, price)
	if err != nil {
		ts.reportError(transNum, "SET_SELL_TRIGGER", user, "No existing sell trigger for this user and stock",
			stock, nil, price.String())
		return "-1"
	}

	err = ts.UserDatabase.RemoveStock(user, stock, trig.GetAmount().IntPart())
	if err != nil {
		ts.reportError(transNum, "SET_SELL_TRIGGER", user, "Could not remove stock from database: "+err.Error(),
			stock, nil, price.String())
		return "-1"
	}

	err = ts.UserDatabase.AddReserveStock(user, stock, trig.GetAmount().IntPart())
	if err != nil {
		ts.reportError(transNum, "SET_SELL_TRIGGER", user, "Could not add stock to reserve: "+err.Error(),
			stock, nil, price.String())
		return "-1"
	}

	go ts.Logger.SystemEvent(ts.Name, transNum, "SET_SELL_TRIGGER", user, stock, nil, price)
	return "1"

}

// CancelSetSell cancels the SET_SELL associated with the given stock and user
// Pre-Conditions: The user must have had a previously set SET_SELL for the given stock
// Post-Conditions:
// 		(a) The set of the user's sell triggers is updated to remove the sell trigger associated with the specified stock
// 		(b) all user account information is reset to the values they would have been if the given SET_SELL command had not been issued
func (ts TransactionServer) CancelSetSell(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]

	trig, err := ts.TriggerClient.CancelSellTrigger(transNum, user, stock)
	if err != nil {
		ts.reportError(transNum, "CANCEL_SET_SELL", user, "No existing sell trigger for this user and stock",
			stock, nil, nil)
		return "-1"
	}

	reserved, err := ts.UserDatabase.GetReserveStock(user, stock)
	if err != nil {
		ts.reportError(transNum, "CANCEL_SET_SELL", user, "Error getting reserved stock from database: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}

	if reserved < trig.GetAmount().IntPart() {
		ts.reportError(transNum, "CANCEL_SET_SELL", user, "Should not have less that a trigger amount in your reserve account",
			stock, nil, nil)
		return "-1"
	}

	err = ts.UserDatabase.RemoveReserveStock(user, stock, trig.GetAmount().IntPart())
	if err != nil {
		ts.reportError(transNum, "CANCEL_SET_SELL", user, "Error removing reserved stock from database: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}

	err = ts.UserDatabase.AddStock(user, stock, trig.GetAmount().IntPart())
	if err != nil {
		ts.reportError(transNum, "CANCEL_SET_SELL", user, "Error adding stock to database: "+err.Error(),
			stock, nil, nil)
		return "-1"
	}

	return "1"
}

// TriggerSuccess listens for incoming successfully executed triggers from the
// triggerserver.
// Params: TRIGGER_SUCCESS,<user>,<stock>,<price>,<amount>,<action>
// t.username, t.stockname, t.price, t.amount, t.action
// Once a successfully completed trigger is received, complete the transaction
// from a user's reserve account to their main account.
func (ts TransactionServer) TriggerSuccess(transNum int, params ...string) string {
	user := params[0]
	stock := params[1]
	price := params[2]
	amount := params[3]
	action := params[4]
	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return "-1"
	}
	priceDec, err := decimal.NewFromString(price)
	if action == "BUY" {
		err = ts.buyExecute(user, stock, amountDec, priceDec)
		if err != nil {
			return "-1"
		}
		return "1"
	} else if action == "SELL" {
		err = ts.sellExecute(user, stock, amountDec, priceDec)
		if err != nil {
			return "-1"
		}
		return "1"
	}
	return "-1"
}

func (ts TransactionServer) reportError(transNum int, command string, user string, errorMsg string, stock interface{}, filename interface{}, funds interface{}) {
	go ts.Logger.SystemError(ts.Name, transNum, command, user, stock, filename, funds,
		errorMsg)
	fmt.Println(errorMsg)
}

func (ts TransactionServer) sellExecute(user string, stock string, amount decimal.Decimal, price decimal.Decimal) error {
	amountShares := amount.IntPart()

	reserved, err := ts.UserDatabase.GetReserveStock(user, stock)
	if err != nil {
		return fmt.Errorf("error getting reserved stock from database:  %s", err.Error())
	}

	if reserved < amountShares {
		return errors.New("reserved stock is less than trigger amount")
	}

	err = ts.UserDatabase.RemoveReserveStock(user, stock, amountShares)
	if err != nil {
		return fmt.Errorf("error removing reserved stock from database:  %s", err.Error())
	}

	err = ts.UserDatabase.AddFunds(user, amount.Mul(price))
	if err != nil {
		return fmt.Errorf("error adding difference between stock cost and reserved:  %s", err.Error())
	}
	return nil
}

func (ts TransactionServer) buyExecute(user string, stock string, amount decimal.Decimal, price decimal.Decimal) error {
	cost, shares, _ := ts.getMaxPurchase(user, stock, amount, price, nil)

	reserved, err := ts.UserDatabase.GetReserveFunds(user)
	if err != nil {
		return fmt.Errorf("error getting reserved funds from database:  %s", err.Error())
	}

	if reserved.LessThan(amount) {
		return errors.New("should not have less than the trigger amount in your reserve account")
	}

	err = ts.UserDatabase.RemoveReserveFunds(user, amount)
	if err != nil {
		return fmt.Errorf("error removing reserved funds: %s", err.Error())
	}

	// Price was lower than the buy trigger
	if amount.GreaterThan(cost) {
		err = ts.UserDatabase.AddFunds(user, amount.Sub(cost))
		if err != nil {
			return fmt.Errorf("error adding difference between stock cost and reserve amount: %s", err.Error())
		}
	}

	err = ts.UserDatabase.AddStock(user, stock, shares)
	if err != nil {
		return fmt.Errorf("error adding stock to database: %s", err.Error())
	}
	return nil
}

// DumpLogUser Print out the history of the users transactions
// to the user specified file
func (ts TransactionServer) DumpLogUser(transNum int, params ...string) string {
	user := params[0]
	filename := params[1]
	go ts.Logger.DumpLog(filename, user)
	return "1"
}

// DisplaySummary provides a summary to the client of the given user's
// transaction history and the current status of their accounts as well
// as any set buy or sell triggers and their parameters.
func (ts TransactionServer) DisplaySummary(transNum int, params ...string) string {
	user := params[0]
	info, err := ts.UserDatabase.GetUserInfo(user)
	if err != nil {
		ts.reportError(transNum, "DISPLAY_SUMMARY", user,
			fmt.Sprintf("Error getting user information from database:  %s", err.Error()), nil, nil, nil)
		return "-1"
	}
	return info
}

// Work with whole numbers for now
// Return the max money you can spend on N shares, given:
// you are user with stock stock and balance balance
func (ts TransactionServer) getMaxPurchase(user string, stock string, availableFunds decimal.Decimal, stockPrice interface{},
	transNum interface{}) (decimal.Decimal, int64, error) {

	var price decimal.Decimal
	if stockPrice != nil {
		price = stockPrice.(decimal.Decimal)
	} else {
		resp, err := quoteclient.Query(user, stock, transNum.(int))
		if err != nil {
			return decimal.Decimal{}, 0, err
		}
		price = resp
	}
	shares := availableFunds.Div(price).IntPart()
	money := price.Mul(decimal.New(shares, 0))
	return money.Round(2), shares, nil
}
