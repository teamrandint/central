package socketserver

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SocketServer struct {
	addr     string
	funcMap  map[string]func(transNum int, args ...string) string
	paramMap map[string]int
	transNum int64
}

func NewSocketServer(addr string) SocketServer {
	return SocketServer{
		addr:     addr,
		funcMap:  make(map[string]func(transNum int, args ...string) string),
		paramMap: make(map[string]int),
		transNum: 0,
	}
}

func (s SocketServer) buildRoutePattern(pattern string) string {
	re := regexp.MustCompile(`(<\w+>)`)
	return re.ReplaceAllString(pattern, `(.+)`) // `(?P\1.+)`
}

func (s SocketServer) Route(key string, f func(transNum int, args ...string) string) {
	s.funcMap[key] = f
}

func (s SocketServer) Run() {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + s.addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		go s.handleRequest(conn)
	}
}

func (s SocketServer) getRoute(command string) (func(transNum int, args ...string) string, []string) {
	command = string(bytes.Trim([]byte(command), "\x00"))
	result := strings.Split(strings.TrimSpace(command), ",")
	function := s.funcMap[result[0]]
	params := result[1:]
	if result[len(result)-1] == "" {
		return nil, nil
	}
	switch result[0] {
	case "COMMIT_BUY", "CANCEL_BUY", "COMMIT_SELL", "CANCEL_SELL", "DISPLAY_SUMMARY":
		if len(params) != 1 {
			return nil, nil
		}
		break
	case "ADD", "QUOTE", "CANCEL_SET_BUY", "CANCEL_SET_SELL":
		if len(params) != 2 {
			return nil, nil
		}
		break
	case "BUY", "SELL", "SET_BUY_AMOUNT", "SET_BUY_TRIGGER", "SET_SELL_TRIGGER", "SET_SELL_AMOUNT":
		if len(params) != 3 {
			return nil, nil
		}
		break
	case "DUMPLOG":
		if len(params) != 1 || len(params) != 2 {
			return nil, nil
		}
	case "TRIGGER_SUCCESS":
		if len(params) != 5 {
			return nil, nil
		}
	default:
		return nil, nil
	}
	return function, params
}

// Handles incoming requests.
func (s SocketServer) handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		recv, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if recv == "\n" { //weird...
			fmt.Println("Received only a newline")
			continue
		}
		fmt.Println("recvd: ", recv)

		sepTransCommand := strings.Split(recv, ";")
		transNum, _ := strconv.Atoi(sepTransCommand[0])
		command := sepTransCommand[1]
		function, params := s.getRoute(command)

		if function == nil {
			fmt.Printf("Error: command not implemented '%s'\n", command)
			n, err := conn.Write([]byte("-1"))
			if err != nil {
				fmt.Println("ERR writing back response ", err)
			} else {
				fmt.Println("Wrote back ", n, " bytes")
			}
			return
		}

		res := function(transNum, params...)
		res += "\n"
		fmt.Println(res)

		// Send a response back to person contacting us.
		n, err := conn.Write([]byte(res))
		if err != nil {
			fmt.Println("ERR writing back response ", err)
		} else {
			fmt.Println("Wrote back ", n, " bytes")
		}
	}
}
