package triggerclient

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type trigger struct {
	username  string
	stockname string
	amount    decimal.Decimal
	price     decimal.Decimal
	action    string
	transNum  int
}

func (t trigger) getPriceStr() string {
	return t.price.String()
}

func (t trigger) getAmountStr() string {
	return t.amount.String()
}

func (t trigger) String() string {
	str := fmt.Sprintf("{%v %v %v %v %v}", t.username, t.stockname, t.getPriceStr(), t.amount, t.action)
	return str
}

func (t trigger) GetCost() decimal.Decimal {
	return t.amount.Mul(t.price)
}

func (t trigger) GetAmount() decimal.Decimal {
	return t.amount
}

func (t trigger) GetPrice() decimal.Decimal {
	return t.price
}

func newSellTrigger(transNum int, username string, stockname string, amount decimal.Decimal) trigger {
	t := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stockname,
		amount:    amount,
		action:    "SELL",
	}

	return t
}

func newBuyTrigger(transNum int, username string, stockname string, amount decimal.Decimal) trigger {
	t := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stockname,
		amount:    amount,
		action:    "BUY",
	}

	return t
}
