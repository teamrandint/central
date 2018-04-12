package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"seng468/quoteserver/logger"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	// _ "net/http/pprof"
)

type QuoteReply struct {
	quote decimal.Decimal
	stock string
	user  string
	time  uint64
	key   string
}

func getReply(msg string) *QuoteReply {
	params := strings.Split(msg, ",")
	if len(params) > 5 {
		return nil
	}

	quote, err := decimal.NewFromString(params[0])
	if err != nil {
		return nil
	}
	timestamp, err := strconv.ParseUint(params[3], 10, 64)
	if err != nil {
		return nil
	}
	return &QuoteReply{
		quote: quote,
		stock: params[1],
		user:  params[2],
		time:  timestamp,
		key:   params[4],
	}
}

func quote(user string, stock string, transNum int) (decimal.Decimal, error) {
	quote, found := quoteCache.Get(stock)
	if found {
		d, _ := decimal.NewFromString(quote.(string))
		return d, nil
	}

	var conn net.Conn
	var err error
	for {
		conn, err = net.DialTimeout("tcp",
			os.Getenv("legacyquoteaddr")+":"+os.Getenv("legacyquoteport"),
			time.Second*1,
		)
		if err != nil { // trans server down? retry
			fmt.Println(err.Error())
		} else {
			break
		}
	}

	request := fmt.Sprintf("%s,%s\n", stock, user)

	conn.Write([]byte(request))
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return decimal.Decimal{}, err
	}
	defer conn.Close()
	reply := getReply(message)
	if reply == nil {
		return decimal.Decimal{}, errors.New("reply from quoteserve doesn't match regex")
	}
	go auditServer.QuoteServer("quoteserver", transNum, reply.quote.String(), reply.stock,
		reply.user, reply.time, reply.key)
	quoteCache.Set(reply.stock, reply.quote.String(), cache.DefaultExpiration)
	return reply.quote, nil
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	user := strings.TrimSpace(query.Get("user"))
	stock := strings.TrimSpace(query.Get("stock"))
	transNum, _ := strconv.Atoi(query.Get("transNum"))
	reply, err := quote(user, stock, transNum)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error receiving quote from legacy quote server", err)
		return
	}
	fmt.Fprintf(w, reply.StringFixed(2))
}

var quoteCache = cache.New(time.Minute, time.Minute)
var auditServer = logger.AuditLogger{Addr: "http://" + os.Getenv("auditaddr") + ":" + os.Getenv("auditport")}

func main() {
	http.HandleFunc("/quote", quoteHandler)
	addr := os.Getenv("quoteaddr")
	port := os.Getenv("quoteport")
	fmt.Printf("Quote server listening on %s:%s\n", addr, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
