package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("Launching server...")
	http.HandleFunc("/userCommand", setBuyTriggerHandler)
	http.HandleFunc("/quoteServer", cancelSetBuyHandler)
	http.HandleFunc("/accountTransaction", setSellTriggerHandler)
	http.HandleFunc("/systemEvent", cancelSetSellHandler)

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

func setBuyTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func cancelSetBuyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func setSellTriggerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func cancelSetSellHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
