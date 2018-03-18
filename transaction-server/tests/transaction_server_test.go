package tests

import (
	"testing"

	"github.com/shopspring/decimal"
)

func NewMockTransactionServer() transactionserver.TransactionServer {
	mockQuote := NewMockQuoteClient()
	mockDB := NewMockDatabase()
	mockLogger := MockLogger{}
	mockServer := MockServer{}
	return transactionserver.TransactionServer{
		Name:         "mock_transaction_serve",
		Addr:         "mock_addr",
		Server:       mockServer,
		Logger:       mockLogger,
		UserDatabase: mockDB,
		QuoteClient:  mockQuote,
	}
}

func TestTransactionServer_Add(t *testing.T) {
	ts := NewMockTransactionServer()
	ts.Add("user1", "50.00")
	actual, _ := ts.UserDatabase.GetFunds("user1")
	expected := decimal.NewFromFloat(50.00)
	if !actual.Equal(expected) {
		t.Error("UserDatabase did not add funds")
	}
}
