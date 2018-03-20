package triggerclient

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

type trigger struct {
	username  string
	stockname string
	amount    int
	price     decimal.Decimal
	action    string
	transNum  int
}

func (t trigger) getPriceStr() string {
	return t.price.String()
}

func (t trigger) getAmountStr() string {
	return strconv.Itoa(t.amount)
}

func (t trigger) String() string {
	str := fmt.Sprintf("{%v %v %v %v %v}", t.username, t.stockname, t.getPriceStr(), t.amount, t.action)
	return str
}

func newSellTrigger(transNum int, username string, stockname string, amount int) trigger {
	t := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stockname,
		amount:    amount,
		action:    "SELL",
	}

	return t
}

func newBuyTrigger(transNum int, username string, stockname string, amount int) trigger {
	t := trigger{
		transNum:  transNum,
		username:  username,
		stockname: stockname,
		amount:    amount,
		action:    "BUY",
	}

	return t
}
