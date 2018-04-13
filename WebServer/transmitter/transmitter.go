package transmitter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/fatih/pool"
)

type Transmitters interface {
	MakeRequest() string
}

type Transmitter struct {
	address        string
	port           string
	connection     net.Conn
	connectionPool pool.Pool
}

func NewTransmitter(addr string, prt string) *Transmitter {
	transmitter := new(Transmitter)
	transmitter.address = addr
	transmitter.port = prt
	factory := func() (net.Conn, error) {
		return net.DialTimeout(
			"tcp",
			addr+":"+prt,
			time.Second*5,
		)
	}
	var err error
	transmitter.connectionPool, err = pool.NewChannelPool(100, 1500, factory)

	// This is real bad and should abort the entire webserver
	if err != nil {
		panic(err)
	}

	return transmitter
}

func (trans *Transmitter) MakeRequest(transNum int, message string) string {
	prefix := strconv.Itoa(transNum)
	message = prefix + ";" + message
	message += "\n"

	conn, err := trans.connectionPool.Get()
	if err != nil {
		fmt.Println("ERROR1: ", err.Error())
	}
	defer conn.Close()

	n, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("ERROR2: ", err)
		pc, _ := conn.(*pool.PoolConn)
		pc.MarkUnusable()
		pc.Close()
	} else {
		fmt.Println("wrote ", n, " bytes")
	}

	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("ERROR3: ", err)
	}
	fmt.Println("recvd: ", reply)

	return reply
}

func (trans *Transmitter) RetrieveDumplog(filename string) []byte {
	auditAddr := "http://" + os.Getenv("auditaddr") + ":" + os.Getenv("auditport")
	resp, err := http.PostForm(auditAddr+"/dumpLogRetrieve", url.Values{"filename": {filename}})
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return body
}
