package tests

import (
	"github.com/shopspring/decimal"
)

type MockDatabase struct {
	userFunds map[string]decimal.Decimal
}

func (db MockDatabase) GetStock(user string, stock string) (int, error) {
	panic("implement me")
}

func (db MockDatabase) RemoveStock(user string, stock string, amount int) error {
	panic("implement me")
}

func (db MockDatabase) PushSell(user string, stock string, cost decimal.Decimal, shares int) error {
	panic("implement me")
}

func (db MockDatabase) PopSell(user string) (stock string, cost decimal.Decimal, shares int, err error) {
	panic("implement me")
}

func (db MockDatabase) AddBuyTrigger(user string, stock string, trigger *transactionserver.Trigger) error {
	panic("implement me")
}

func (db MockDatabase) GetBuyTrigger(user string, stock string) (*transactionserver.Trigger, error) {
	panic("implement me")
}

func (db MockDatabase) RemoveBuyTrigger(user string, stock string) (*transactionserver.Trigger, error) {
	panic("implement me")
}

func (db MockDatabase) AddSellTrigger(user string, stock string, trigger *transactionserver.Trigger) error {
	panic("implement me")
}

func (db MockDatabase) GetSellTrigger(user string, stock string) (*transactionserver.Trigger, error) {
	panic("implement me")
}

func (db MockDatabase) RemoveSellTrigger(user string, stock string) (*transactionserver.Trigger, error) {
	panic("implement me")
}

func NewMockDatabase() MockDatabase {
	return MockDatabase{
		userFunds: make(map[string]decimal.Decimal),
	}
}

func (db MockDatabase) AddFunds(user string, amount decimal.Decimal) error {
	if val, ok := db.userFunds[user]; ok {
		db.userFunds[user] = val.Add(amount)
		return nil
	}
	db.userFunds[user] = amount
	return nil
}

func (db MockDatabase) GetFunds(user string) (decimal.Decimal, error) {
	if val, ok := db.userFunds[user]; ok {
		return val, nil
	}
	return decimal.NewFromFloat(0), nil
}

func (db MockDatabase) RemoveFunds(string, decimal.Decimal) error {
	panic("implement me")
}

func (db MockDatabase) PushBuy(user string, stock string, cost decimal.Decimal, shares int) error {
	panic("implement me")
}

func (db MockDatabase) PopBuy(user string) (stock string, cost decimal.Decimal, shares int, err error) {
	panic("implement me")
}

func (db MockDatabase) AddStock(user string, stock string, shares int) error {
	panic("implement me")
}
