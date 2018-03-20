package triggerclient

import (
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/shopspring/decimal"
)

var triggeraddr = os.Getenv("triggeraddr")
var triggerport = os.Getenv("triggerport")
var triggerURL = "http://" + triggeraddr + ":" + triggerport

const (
	setEndpoint    = "/setTrigger"
	startEndpoint  = "/startTrigger"
	cancelEndpoint = "/cancelTrigger"
	listEndpoint   = "/runningTriggers"
)

// SetNewSellTrigger adds a new sell trigger to the triggerserver
func SetNewSellTrigger(transNum int, username string, stock string, amount int) error {
	trig := newSellTrigger(transNum, username, stock, amount)
	return setTrigger(transNum, trig)
}

// SetSellTrigger adds a new sell trigger to the triggerserver
func SetSellTrigger(transNum int, trig trigger) error {
	return setTrigger(transNum, trig)
}

// StartSellTrigger adds a new sell trigger to the triggerserver
func StartSellTrigger(transNum int, trig trigger) error {
	return startTrigger(transNum, trig)
}

// StartNewSellTrigger starts an existing sell trigger on the triggerserver
func StartNewSellTrigger(transNum int, username string, stock string, price decimal.Decimal) error {
	trig := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stock,
		price:     price,
		action:    "SELL",
	}
	return startTrigger(transNum, trig)
}

// SetNewBuyTrigger adds a new sell trigger to the triggerserver
func SetNewBuyTrigger(transNum int, username string, stock string, amount int) error {
	trig := newBuyTrigger(transNum, username, stock, amount)
	return setTrigger(transNum, trig)
}

// SetBuyTrigger adds a new Buy trigger to the triggerserver
func SetBuyTrigger(transNum int, trig trigger) error {
	return setTrigger(transNum, trig)
}

// StartBuyTrigger adds a new buy trigger to the triggerserver
func StartBuyTrigger(transNum int, trig trigger) error {
	return startTrigger(transNum, trig)
}

// StartNewBuyTrigger starts an existing Buy trigger on the triggerserver
func StartNewBuyTrigger(transNum int, username string, stock string, price decimal.Decimal) error {
	trig := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stock,
		price:     price,
		action:    "BUY",
	}
	return startTrigger(transNum, trig)
}

// setTrigger adds a new trigger to the triggerserver.
// Action is either 'BUY' or 'SELL'
func setTrigger(transNum int, newTrigger trigger) error {
	values := url.Values{
		"action":   {newTrigger.action},
		"transnum": {strconv.Itoa(transNum)},
		"username": {newTrigger.username},
		"stock":    {newTrigger.stockname},
		"amount":   {newTrigger.getAmountStr()},
	}
	_, err := http.PostForm(triggerURL+startEndpoint, values)
	if err != nil {
		return err
	}

	return nil
}

// startTrigger starts an existing trigger on the triggerserver.
func startTrigger(transNum int, newTrigger trigger) error {
	values := url.Values{
		"action":   {newTrigger.action},
		"transnum": {strconv.Itoa(transNum)},
		"username": {newTrigger.username},
		"stock":    {newTrigger.stockname},
		"price":    {newTrigger.getPriceStr()},
	}
	_, err := http.PostForm(triggerURL+setEndpoint, values)
	if err != nil {
		return err
	}

	return nil
}

// CancelTrigger cancels a running trigger on the triggerserver.
// Action is either 'BUY' or 'SELL'
// If the given trigger could not be found, returns an error.
func CancelTrigger(transNum int, cancel trigger) error {
	values := url.Values{
		"action":   {cancel.action},
		"transnum": {strconv.Itoa(transNum)},
		"username": {cancel.username},
		"stock":    {cancel.stockname},
	}
	_, err := http.PostForm(triggerURL+cancelEndpoint, values)
	if err != nil {
		return err
	}

	return nil
}

// ListRunningTriggers returns a list of all running triggers on the TriggerServer
// TODO: something useful if needed
func ListRunningTriggers() {
	_, err := http.Get(triggerURL + listEndpoint)
	if err != nil {
		panic(err)
	}
}
