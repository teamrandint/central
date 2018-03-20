package main

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var triggeraddr = os.Getenv("triggeraddr")
var triggerport = os.Getenv("triggerport")
var triggerURL = "http://" + triggeraddr + ":" + triggerport

const (
	setEndpoint    = "/setTrigger"
	cancelEndpoint = "/cancelTrigger"
	listEndpoint   = "/runningTriggers"
)

// SetSellTrigger adds a new sell trigger to the triggerserver
func SetSellTrigger(transNum int, trig trigger) error {
	return setTrigger(transNum, "SELL", trig)
}

// SetBuyTrigger adds a new Buy trigger to the triggerserver
func SetBuyTrigger(transNum int, trig trigger) error {
	return setTrigger(transNum, "BUY", trig)
}

// setTrigger adds a new trigger to the triggerserver.
// Action is either 'BUY' or 'SELL'
// Once a trigger is added, it automatically runs
func setTrigger(transNum int, action string, newTrigger trigger) error {
	values := url.Values{
		"action":   {action},
		"transnum": {strconv.Itoa(transNum)},
		"username": {newTrigger.username},
		"stock":    {newTrigger.stockname},
		"price":    {newTrigger.getPriceStr()},
	}
	_, err := http.PostForm(triggerURL+setEndpoint, values)
	if err != nil {
		return err
	}

	return nil
}

// CancelTrigger cancels a running trigger on the triggerserver.
// Action is either 'BUY' or 'SELL'
// If the given trigger could not be found, returns an error.
func CancelTrigger(transNum int, action string, cancel trigger) error {
	values := url.Values{
		"action":   {action},
		"transnum": {strconv.Itoa(transNum)},
		"username": {cancel.username},
		"stock":    {cancel.stockname},
		"price":    {cancel.getPriceStr()},
	}
	_, err := http.PostForm(triggerURL+cancelEndpoint, values)
	if err != nil {
		return err
	}

	return nil
}

// ListRunningTriggers returns a list of all running triggers on the TriggerServer
func ListRunningTriggers() {
	_, err := http.Get(triggerURL + listEndpoint)
	if err != nil {
		panic(err)
	}
}
