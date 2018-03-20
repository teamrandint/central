package triggerclient

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

// CancelSellTrigger attempts to cancel an existing sell trigger on the server
func CancelSellTrigger(transNum int, username string, stock string) (trigger, error) {
	trig := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stock,
		action:    "SELL",
	}
	return cancelTrigger(transNum, trig)
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

// Attempts to cancel an existing Buy trigger on the server
func CancelBuyTrigger(transNum int, username string, stock string) (trigger, error) {
	trig := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stock,
		action:    "BUY",
	}
	return cancelTrigger(transNum, trig)
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
// Should return the cancelled triggers details
func cancelTrigger(transNum int, cancel trigger) (trigger, error) {
	values := url.Values{
		"action":   {cancel.action},
		"transnum": {strconv.Itoa(transNum)},
		"username": {cancel.username},
		"stock":    {cancel.stockname},
	}
	resp, err := http.PostForm(triggerURL+cancelEndpoint, values)
	if err != nil {
		return trigger{}, err
	}

	trig := getTriggerFromResponse(resp)
	return trig, nil
}

// ListRunningTriggers returns a list of all running triggers on the TriggerServer
// TODO: something useful if needed
func ListRunningTriggers() {
	_, err := http.Get(triggerURL + listEndpoint)
	if err != nil {
		panic(err)
	}
}

func getTriggerFromResponse(resp *http.Response) trigger {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	return parseTriggerFromString(bodyString)
}

// "{%v %v %v %v %v}", t.username, t.stockname, t.getPriceStr(), t.amount, t.action
func parseTriggerFromString(trigStr string) trigger {
	re := regexp.MustCompile(`{(\w+) (\w+) (\d+.\d+) (\d+) (\w+)}`)
	matches := re.FindStringSubmatch(trigStr)

	price, err := decimal.NewFromString(matches[3])
	if err != nil {
		panic(err)
	}
	amount, err := strconv.Atoi(matches[4])
	if err != nil {
		panic(err)
	}

	trig := trigger{
		username:  matches[1],
		stockname: matches[2],
		price:     price,
		amount:    amount,
		action:    matches[5],
	}
	return trig
}
