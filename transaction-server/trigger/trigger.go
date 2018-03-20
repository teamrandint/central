package triggerclient

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type trigger struct {
	username  string
	stockname string
	price     decimal.Decimal
	action    string
	transNum  int
}

func (t trigger) getPriceStr() string {
	return t.price.String()
}

func (t trigger) String() string {
	str := fmt.Sprintf("{%v %v %v %v}", t.username, t.stockname, t.getPriceStr(), t.action)
	return str
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
