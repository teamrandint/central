package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// runningTriggers follows [action][stock][user] indexing
type runningTriggersKey struct {
	action, stock, user string
}

var runningTriggers = make(map[runningTriggersKey]trigger)
var runningTriggersLock sync.RWMutex

var successListener = make(chan trigger)

func main() {
	fmt.Println("Launching server...")
	http.HandleFunc("/setTrigger", setTriggerHandler)
	http.HandleFunc("/cancelTrigger", cancelTriggerHandler)
	http.HandleFunc("/runningTriggers", getRunningTriggersHandler)

	go startSuccessListener()

	fmt.Printf("Trigger server listening on %s:%s\n", os.Getenv("triggeraddr"), os.Getenv("triggerport"))
	if err := http.ListenAndServe(":"+os.Getenv("triggerport"), nil); err != nil {
		panic(err)
	}
}

func setTriggerHandler(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	transnumStr := r.FormValue("transnum")
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

	transnum, err := strconv.Atoi(transnumStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	var t trigger
	if action == "BUY" {
		t = newBuyTrigger(successListener, transnum, username, stock, price)
	} else {
		t = newSellTrigger(successListener, transnum, username, stock, price)
	}
	fmt.Println("Added: ", t)
	go t.StartPolling()

	runningTriggersLock.Lock()
	runningTriggers[runningTriggersKey{t.action, t.stockname, t.username}] = t
	runningTriggersLock.Unlock()

	w.WriteHeader(http.StatusOK)
}

func cancelTriggerHandler(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	transnumStr := r.FormValue("transnum")
	username := r.FormValue("username")
	stock := r.FormValue("stock")

	if !verifyAction(action) {
		w.WriteHeader(http.StatusBadRequest)
		panic("Tried to post a bad action (BUY/SELL)")
	}

	transnum, err := strconv.Atoi(transnumStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	t := findRunningTrigger(transnum, action, username, stock)
	err = cancelTrigger(t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	fmt.Println(t)

	w.WriteHeader(http.StatusOK)
}

func startSuccessListener() {
	select {
	case trig := <-successListener:
		fmt.Println("Closing successful trigger: ", trig)
		go cancelTrigger(trig)
		go alertTriggerSuccess(trig)
	}
}

// Send an alert back to the transaction server when a trigger successfully finishes
func alertTriggerSuccess(t trigger) {
	conn, err := net.DialTimeout("tcp",
		os.Getenv("transaddr")+":"+os.Getenv("transport"),
		time.Second*15,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(t.getSuccessString())
	_, err = fmt.Fprintf(conn, t.getSuccessString())
	if err != nil {
		panic(err)
	}

	err = conn.Close()
	if err != nil {
		panic(err)
	}
}

func getRunningTriggersHandler(w http.ResponseWriter, r *http.Request) {
	runningTriggersLock.RLock()
	fmt.Fprintln(w, runningTriggers)
	runningTriggersLock.RUnlock()
}

func verifyAction(action string) bool {
	if action != "BUY" && action != "SELL" {
		return false
	}
	return true
}

func findRunningTrigger(transnum int, action string, username string, stock string) trigger {
	runningTriggersLock.RLock()
	defer runningTriggersLock.RUnlock()
	return runningTriggers[runningTriggersKey{action, stock, username}]
}

// Removes the trigger from the poller
func cancelTrigger(t trigger) error {
	runningTriggersLock.Lock()
	_, ok := runningTriggers[runningTriggersKey{t.action, t.stockname, t.username}]
	if !ok {
		return errors.New("Can't find running trigger")
	}
	runningTriggers[runningTriggersKey{t.action, t.stockname, t.username}].Cancel()
	delete(runningTriggers, runningTriggersKey{t.action, t.stockname, t.username})
	runningTriggersLock.Unlock()
	return nil
}
