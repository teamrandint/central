package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type trigger struct {
	username  string
	stockname string
	price     decimal.Decimal
	action    string
}

func (t trigger) getSuccessString() string {
	return fmt.Sprintf("TRIGGER_SUCCESS,%v,%v,%v,%v\n",
		t.username, t.stockname, t.price, t.action)
}

func newSellTrigger(username string, stockname string, price decimal.Decimal) trigger {
	t := trigger{
		username:  username,
		stockname: stockname,
		price:     price,
		action:    "SELL",
	}

	return t
}

func newBuyTrigger(username string, stockname string, price decimal.Decimal) trigger {
	t := trigger{
		username:  username,
		stockname: stockname,
		price:     price,
		action:    "BUY",
	}

	return t
}
