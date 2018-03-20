package main

import (
	"fmt"
	"net"
	"os"
	"seng468/triggerserver/quote"
	"time"

	"github.com/shopspring/decimal"
)

type trigger struct {
	username  string
	stockname string
	price     decimal.Decimal
	action    string
	transNum  int
	done      chan bool
}

func (t trigger) getSuccessString() string {
	return fmt.Sprintf("TRIGGER_SUCCESS,%v,%v,%v,%v\n",
		t.username, t.stockname, t.price, t.action)
}

func (t trigger) getPriceStr() string {
	return t.price.String()
}

func (t trigger) String() string {
	str := fmt.Sprintf("{%v %v %v %v}", t.username, t.stockname, t.getPriceStr(), t.action)
	return str
}

func (t trigger) StartPolling() {
	if t.checkTriggerStatus() {
		t.Cancel()
		t.alertTriggerSuccess()
	}
	ticker := time.NewTicker((time.Second * 60) + time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-t.done:
			return
		case <-ticker.C:
			if t.checkTriggerStatus() {
				t.Cancel()
				t.alertTriggerSuccess()
			}
		}
	}
}

// Send an alert back to the transaction server when a trigger successfully finishes
func (t trigger) alertTriggerSuccess() {
	conn, err := net.DialTimeout("tcp",
		os.Getenv("transaddr")+":"+os.Getenv("transport"),
		time.Second*15,
	)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(conn, t.getSuccessString())
	if err != nil {
		panic(err)
	}

	err = conn.Close()
	if err != nil {
		panic(err)
	}
}

func (t trigger) Cancel() {
	t.done <- true
}

func (t trigger) checkTriggerStatus() bool {
	result := t.hitQuoteServer()
	return t.checkResult(result)
}

func (t trigger) hitQuoteServer() decimal.Decimal {
	result, err := quoteclient.Query(t.username, t.stockname, t.transNum) //user string, stock string, transNum int
	if err != nil {
		panic(err)
	}

	return result
}

// See if the result from the quoteserver is enough to stop the trigger
func (t trigger) checkResult(result decimal.Decimal) bool {
	switch t.action {
	case "BUY":
		return t.price.LessThanOrEqual(result)
	case "SELL":
		return t.price.GreaterThanOrEqual(result)
	}

	panic("Should never reach here...")
}

func newSellTrigger(transNum int, username string, stockname string, price decimal.Decimal) trigger {
	t := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stockname,
		price:     price,
		action:    "SELL",
	}

	return t
}

func newBuyTrigger(transNum int, username string, stockname string, price decimal.Decimal) trigger {
	t := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stockname,
		price:     price,
		action:    "BUY",
	}

	return t
}
