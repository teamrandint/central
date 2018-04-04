package triggerclient

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Trigger struct {
	username  string
	stockname string
	amount    int64
	price     decimal.Decimal
	action    string
	transNum  int
}

func (t Trigger) getPriceStr() string {
	return t.price.String()
}

func (t Trigger) getAmountStr() string {
	return string(t.amount)
}

func (t Trigger) String() string {
	str := fmt.Sprintf("{%v %v %v %v %v}", t.username, t.stockname, t.getPriceStr(), t.getAmountStr(), t.action)
	return str
}

func (t Trigger) GetCost() decimal.Decimal {
	return t.price.Mul(decimal.NewFromFloat(float64(t.amount)))
}

func (t Trigger) GetAmount() int64 {
	return t.amount
}

func (t Trigger) GetPrice() decimal.Decimal {
	return t.price
}

func newSellTrigger(transNum int, username string, stockname string, amount int64) Trigger {
	t := Trigger{
		transNum:  transNum,
		username:  username,
		stockname: stockname,
		amount:    amount,
		action:    "SELL",
	}

	return t
}

func newBuyTrigger(transNum int, username string, stockname string, amount int64) Trigger {
	t := Trigger{
		transNum:  transNum,
		username:  username,
		stockname: stockname,
		amount:    amount,
		action:    "BUY",
	}

	return t
}
