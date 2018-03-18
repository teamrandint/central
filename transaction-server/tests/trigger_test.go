package tests

import (
	"seng468/transaction-server/trigger"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func Callback(t *testing.T, expected *triggers.Trigger, called *bool) func(trigger *transactionserver.Trigger) {
	return func(trigger *triggers.Trigger) {
		if expected.User != trigger.User {
			t.Error("User name does not match")
		}
		if expected.Stock != trigger.Stock {
			t.Error("Stock does not match")
		}
		if expected.TransNum != trigger.TransNum {
			t.Error("Transaction number does not match")
		}
		if expected.QuoteClient != trigger.QuoteClient {
			t.Error("Quote client does not match")
		}
		if expected.BuySellAmount != trigger.BuySellAmount {
			t.Error("Buy amount does not match")
		}
		if expected.TriggerAmount != trigger.TriggerAmount {
			t.Error("Trigger amount does not match")
		}
		if expected.TriggerType != trigger.TriggerType {
			t.Error("Trigger amount does not match")
		}
		*called = true
	}
}

func TestTrigger_Buy(t *testing.T) {
	mockQuote := NewMockQuoteClient()
	mockQuote.addRule("ABC", decimal.NewFromFloat(21.00))
	buyAmount := decimal.NewFromFloat(10.00)
	triggerAmount := decimal.NewFromFloat(20.00)
	expected := &transactionserver.Trigger{
		User:          "user",
		Stock:         "ABC",
		TransNum:      1,
		QuoteClient:   mockQuote,
		BuySellAmount: buyAmount,
		TriggerAmount: triggerAmount,
		TriggerType:   "BUY",
	}
	called := false
	trig := transactionserver.NewBuyTrigger("user", "ABC", mockQuote, buyAmount, Callback(t, expected, &called))
	time.Sleep(time.Second)
	trig.Start(triggerAmount, 1)
	time.Sleep(time.Second)
	if called {
		t.Error("Trigger called too early")
	}
	mockQuote.addRule("ABC", decimal.NewFromFloat(19.00))
	time.Sleep(time.Second * 3)
	if !called {
		t.Error("Trigger was never called")
	}
}
