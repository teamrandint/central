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
	address    string
	port       string
	connection net.Conn
	connectionPool pool.Pool
}

func NewTransmitter(addr string, prt string) *Transmitter {
	transmitter := new(Transmitter)
	transmitter.address = addr
	transmitter.port = prt
	factory := func() (net.Conn, error) { return net.Dial("tcp", trans.address+":"+trans.port) }

	transmitter.connectionPool.NewChannelPool(5, 30, factory)

	return transmitter
}

func (trans *Transmitter) MakeRequest(transNum int, message string) string {
	prefix := strconv.Itoa(transNum)
	message = prefix + ";" + message
	message += "\n"
	var conn net.Conn
	var err error
	for {
		conn, err = net.DialTimeout(
			"tcp",
			trans.address+":"+trans.port,
			time.Second*5,
		)

		if err != nil { // trans server down? retry
			fmt.Println("Trans server timedout -- retrying")
		} else {
			break
		}
	}

	trans.connection = conn

	trans.connection.Write([]byte(message))
	reply, _ := bufio.NewReader(trans.connection).ReadString('\n')
	trans.connection.Close()
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
