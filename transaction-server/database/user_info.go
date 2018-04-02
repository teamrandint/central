package database

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"github.com/shopspring/decimal"
)

type UserInfo struct {
	user string
	funds float64
	reservedFunds float64
	reservedStock map[string]string
	stock map[string]string
	sellOrders []string
	buyOrders []string
}

func GetUserInfoFromReply(user string, reply interface{}) (UserInfo, error) {
	var balance float64
	var stock interface{}
	var sellOrders []string
	var buyOrders []string
	var reservedFunds float64
	var stockReserve interface{}
	var stockMap map[string]string
	var stockReserveMap map[string]string

	values, err := redis.Values(reply, nil); if err != nil {
		return UserInfo{}, err
	}

	if _, err := redis.Scan(values, &balance, &stock, &sellOrders, &buyOrders, &reservedFunds, &stockReserve); err != nil {
		return UserInfo{}, err
	}

	stockReserveMap, err = redis.StringMap(stockReserve, err); if err != nil {
		return UserInfo{}, err
	}

	stockReserveMap, err = redis.StringMap(stock, err); if err != nil {
		return UserInfo{}, err
	}

	return UserInfo{
		user: user,
		funds: balance,
		reservedFunds: reservedFunds,
		reservedStock: stockReserveMap,
		stock: stockMap,
		sellOrders: sellOrders,
		buyOrders: buyOrders,
	}, nil
}

func (info UserInfo) getString() string {
	str := fmt.Sprintf("User:\t\t\t%s;Funds:\t\t\t%.2f;", info.user, info.funds)
	if len(info.stock) > 0 {
		str += "Stock:;"
	}
	for stock, amount := range info.stock {
		dec, _ := decimal.NewFromString(amount); if dec.GreaterThan(decimal.Zero) {
			str += fmt.Sprintf("\t%s:\t%s\n", stock, amount)
		}
	}

	if len(info.buyOrders) > 0 {
		str += "Buy Orders:;"
	}
	for _, buyOrder := range info.buyOrders {
		stock, cost, _ := decodeOrder(buyOrder)
		if cost.GreaterThan(decimal.Zero) {
			str += fmt.Sprintf("\t%s:\t%s;", stock, cost.StringFixed(2))
		}
	}

	if len(info.sellOrders) > 0 {
		str += "Sell Orders:;"
	}
	for _, sellOrder := range info.sellOrders {
		stock, cost, _ := decodeOrder(sellOrder)
		if cost.GreaterThan(decimal.Zero) {
			str += fmt.Sprintf("\t%s:\t%s;", stock, cost.StringFixed(2))
		}
	}

	str += fmt.Sprintf("Reserved Funds:\t%.2f;", info.reservedFunds)

	if len(info.reservedStock) > 0 {
		str += "Reserved stock:;"
	}
	for key, value := range info.reservedStock {
		dec, _ := decimal.NewFromString(value); if dec.GreaterThan(decimal.Zero) {
			str += fmt.Sprintf("\t%s:\t%s;", key, value)
		}
	}
	str += "\n"
	return str
}