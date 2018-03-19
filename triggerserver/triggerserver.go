package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/shopspring/decimal"
)

func main() {
	fmt.Println("Launching server...")
	http.HandleFunc("/setTrigger", setTriggerHandler)
	http.HandleFunc("/cancelTrigger", cancelTriggerHandler)
	http.HandleFunc("/runningTriggers", getRunningTriggersHandler)

	fmt.Printf("Trigger server listening on %s:%s\n", os.Getenv("triggeraddr"), os.Getenv("triggerport"))
	if err := http.ListenAndServe(":"+os.Getenv("triggerport"), nil); err != nil {
		panic(err)
	}
}

// Send an alert back to the transaction server when a trigger successfully finishes
func alertTriggerSuccess(finished trigger) {
	conn, err := net.DialTimeout("tcp",
		os.Getenv("transaddr")+":"+os.Getenv("transport"),
		time.Second,
	)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(conn, finished.getSuccessString())
	if err != nil {
		panic(err)
	}

	err = conn.Close()
	if err != nil {
		panic(err)
	}
}

func setTriggerHandler(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	transnum := r.FormValue("transnum")
	username := r.FormValue("username")
	stock := r.FormValue("stock")
	priceStr := r.FormValue("price")

	if !verifyAction(action) {
		w.WriteHeader(http.StatusBadRequest)
		panic("Tried to post a bad action (BUY/SELL)")
	}

	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	var t trigger
	if action == "BUY" {
		t = newBuyTrigger(username, stock, price)
	} else {
		t = newSellTrigger(username, stock, price)
	}
	fmt.Println(transnum, t)

	w.WriteHeader(http.StatusOK)
}

func cancelTriggerHandler(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	transnum := r.FormValue("transnum")
	username := r.FormValue("username")
	stock := r.FormValue("stock")

	if !verifyAction(action) {
		w.WriteHeader(http.StatusBadRequest)
		panic("Tried to post a bad action (BUY/SELL)")
	}

	t := findRunningTrigger(action, username, stock)
	err := cancelTrigger(t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	fmt.Println(transnum, t)

	w.WriteHeader(http.StatusOK)
}

func getRunningTriggersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func verifyAction(action string) bool {
	if action != "BUY" && action != "SELL" {
		return false
	}
	return true
}

func findRunningTrigger(action string, username string, stock string) trigger {
	dec, _ := decimal.NewFromString("11.11")
	return newBuyTrigger(username, stock, dec)
}

func cancelTrigger(t trigger) error {
	return nil
}
