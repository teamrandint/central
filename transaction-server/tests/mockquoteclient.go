package tests

import (
	"errors"

	"github.com/shopspring/decimal"
)

type MockQuoteClient struct {
	stockMap map[string]decimal.Decimal
}

func (qc *MockQuoteClient) Query(user string, stock string, transNum int) (decimal.Decimal, error) {
	if val, ok := qc.stockMap[stock]; ok {
		return val, nil
	}
	return decimal.Decimal{}, errors.New("stock not mocked")
}

func NewMockQuoteClient() *MockQuoteClient {
	return &MockQuoteClient{
		stockMap: make(map[string]decimal.Decimal),
	}
}

func (qc *MockQuoteClient) addRule(stock string, amount decimal.Decimal) {
	qc.stockMap[stock] = amount
}
