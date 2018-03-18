package triggers

import (
	"fmt"
	"seng468/transaction-server/quote"
	"time"

	"github.com/shopspring/decimal"
)

type Trigger struct {
	User          string
	Stock         string
	TransNum      int
	BuySellAmount decimal.Decimal
	TriggerAmount decimal.Decimal
	action        func(trig *Trigger)
	TriggerType   string
	cancel        chan bool
}

func NewBuyTrigger(user string, stock string, buySellAmount decimal.Decimal, action func(*Trigger)) *Trigger {
	return &Trigger{
		User:          user,
		Stock:         stock,
		BuySellAmount: buySellAmount,
		action:        action,
		TriggerType:   "BUY",
	}
}

func NewSellTrigger(user string, stock string, buySellAmount decimal.Decimal, action func(trigger *Trigger)) *Trigger {
	return &Trigger{
		User:          user,
		Stock:         stock,
		BuySellAmount: buySellAmount,
		action:        action,
		TriggerType:   "SELL",
	}
}

func (trig Trigger) Start(trigger decimal.Decimal, transNum int) {
	trig.TriggerAmount = trigger
	trig.TransNum = transNum
	trig.cancel = make(chan bool)
	go func() {
		for {
			trig.testTrigger()
			select {
			case <-time.After(time.Millisecond * 200):
			case <-trig.cancel:
				return
			}
		}
	}()
}

func (trig Trigger) Cancel() {
	if trig.cancel != nil {
		trig.cancel <- true
	}
}

func (trig Trigger) testTrigger() {
	quote, err := quoteclient.Query(trig.User, trig.Stock, trig.TransNum)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if trig.TriggerType == "BUY" && quote.LessThanOrEqual(trig.TriggerAmount) {
		trig.action(&trig)
		trig.cancel <- true
		return
	}

	if trig.TriggerType == "SELL" && quote.GreaterThanOrEqual(trig.TriggerAmount) {
		trig.action(&trig)
		trig.cancel <- true
	}
}
