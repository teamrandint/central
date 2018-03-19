package transmitter

import (
	"bufio"
	"fmt"
	"log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Transmitters interface {
	MakeRequest() string
}

type Transmitter struct {
	address    string
	port       string
	connection net.Conn
}

func NewTransmitter(addr string, prt string) *Transmitter {
	transmitter := new(Transmitter)
	transmitter.address = addr
	transmitter.port = prt

	// Create a connection to the specified server
	var conn net.Conn
	var err error
	for err != nil {
		conn, err = net.Dial("tcp", addr+":"+prt)
		time.Sleep(time.Millisecond * 30)
		log.Print(err)
	}

	transmitter.connection = conn

	return transmitter
}

func (trans *Transmitter) MakeRequest(transNum int, message string) string {
	prefix := strconv.Itoa(transNum)
	message = prefix + ";" + message
	message += "\n"
	conn, err := net.Dial("tcp", trans.address+":"+trans.port)

	if err != nil {
		// Error in connection
		log.Print(err)
		return "-1"
	}
		
	trans.connection = conn

	// fmt.Println("Making request to transaction server")
	fmt.Fprintf(trans.connection, message)
	// fmt.Println("Waiting for response from transaction server")
	reply, _ := bufio.NewReader(trans.connection).ReadString('\n')
	trans.connection.Close()
	return reply
}

func (trans *Transmitter) RetrieveDumplog(filename string) []byte {
	// auditaddr = os.Getenv("auditaddr")
	// auditport = os.Getenv("auditport")
	resp, err := http.PostForm("http://172.20.0.3:44455/dumpLogRetrieve", url.Values{"filename": {filename}})
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
