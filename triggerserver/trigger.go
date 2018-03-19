package main

import (
	"fmt"
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
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()

	for {
		select {
		case <-t.done:
			return
		case <-ticker.C:
			t.hitQuoteServer()
		}
	}
}

func (t trigger) Cancel() {
	t.done <- true
}

func (t trigger) hitQuoteServer() decimal.Decimal {
	result, err := quoteclient.Query(t.username, t.stockname, t.transNum) //user string, stock string, transNum int
	if err != nil {
		panic(err)
	}
	return result
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
