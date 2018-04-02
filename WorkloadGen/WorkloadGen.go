package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var transcount uint64
var endpointTimes map[string][]time.Duration
var endpointMutex sync.Mutex

// go run WorkloadGen.go serverAddr:port workloadfile
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: server address, workloadfile, delay(ms)")
		return
	}

	serverAddr := os.Args[1]
	workloadFile := os.Args[2]
	delayMs, _ := strconv.Atoi(os.Args[3])
	endpointTimes = make(map[string][]time.Duration)

	fmt.Printf("Testing %v on serverAddr %v with delay of %vms\n", workloadFile, serverAddr, delayMs)

	users := splitUsersFromFile(workloadFile)
	fmt.Printf("Found %d users...\n", len(users))
	go countTPS()

	runRequests(serverAddr, users, delayMs)
	fmt.Printf("Done!\n")
}

func runRequests(serverAddr string, users map[string][]string, delay int) {
	var wg sync.WaitGroup
	for userName, commands := range users {
		fmt.Printf("Running user %v's commands...\n", userName)

		wg.Add(1)
		go func(commands []string) {
			//timeout := time.Duration(15 * time.Second)
			client := http.Client{
			//	Timeout: timeout,
			}

			// Issue login before executing any commands
			resp, err := client.PostForm("http://"+serverAddr+"/"+"LOGIN"+"/", url.Values{"username": {userName}})
			if err != nil {
				fmt.Println(err)
			} else {
				resp.Body.Close()
			}

			for _, command := range commands {
				endpoint, values := parseCommand(command)
				time.Sleep(time.Duration(delay) * time.Millisecond) // ADJUST THIS TO CHANGE DELAY
				// fmt.Println("http://"+serverAddr+"/"+endpoint+"/", values)

				var resp *http.Response
				var err error
				var time0 time.Time
				for {
					time0 = time.Now()
					resp, err = client.PostForm("http://"+serverAddr+"/"+endpoint+"/", values)
					if err != nil {
						fmt.Println("Post timed out -- retrying")
					} else {
						break
					}
				}

				resp.Body.Close()
				responseTime := time.Since(time0)
				endpointMutex.Lock()
				endpointTimes[endpoint] = append(endpointTimes[endpoint], responseTime)
				endpointMutex.Unlock()
				atomic.AddUint64(&transcount, 1)
			}

			wg.Done()
		}(commands)
	}

	// Wait for commands, then manually post the final dumplog
	wg.Wait()
	resp, httpErr := http.PostForm("http://"+serverAddr+"/DUMPLOG/", url.Values{"filename": {"./output.xml"}})
	if httpErr != nil {
		panic(httpErr)
	}

	data, decodeErr := ioutil.ReadAll(resp.Body)
	if decodeErr != nil {
		panic(decodeErr)
	}

	file, createErr := os.Create("./output.xml")
	defer file.Close()
	if createErr != nil {
		fmt.Printf("error: %v %v\n", createErr, file)
	}
	_, writeErr := file.Write(data)
	if writeErr != nil {
		panic(writeErr)
	}

	printEndpointStats()
	saveEndpointStats()
}

func splitUsersFromFile(filename string) map[string][]string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	// https://regex101.com/r/O6xaTp/3
	re := regexp.MustCompile(`\[\d+\] ((?P<endpoint>\w+),(?P<user>\w+)(,-*\w*\.*\d*)*)`)
	outputCommands := make(map[string][]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)

		if matches != nil {
			command := matches[1]
			//endpoint := matches[2]
			user := matches[3]
			outputCommands[user] = append(outputCommands[user], command)
		} else {
			fmt.Println("Error parsing command: ", line)
		}
	}

	return outputCommands
}

// Parse a single line command into the corresponding endpoint and values
func parseCommand(cmd string) (endpoint string, v url.Values) {
	subcmd := strings.Split(cmd, ",")
	endpoint = subcmd[0]
	// username, stock, amount, filename
	switch endpoint {
	case "ADD":
		v = url.Values{
			"username": {subcmd[1]},
			"amount":   {subcmd[2]},
		}
	case "QUOTE", "CANCEL_SET_BUY", "CANCEL_SET_SELL":
		v = url.Values{
			"username": {subcmd[1]},
			"stock":    {subcmd[2]},
		}
	case "SELL", "BUY", "SET_BUY_AMOUNT", "SET_BUY_TRIGGER", "SET_SELL_AMOUNT", "SET_SELL_TRIGGER":
		v = url.Values{
			"username": {subcmd[1]},
			"stock":    {subcmd[2]},
			"amount":   {subcmd[3]},
		}
	case "COMMIT_BUY", "CANCEL_BUY", "COMMIT_SELL", "CANCEL_SELL", "DISPLAY_SUMMARY":
		v = url.Values{
			"username": {subcmd[1]},
		}
	}

	return endpoint, v
}

func countTPS() {
	var tpsStart uint64
	var tpsEnd uint64
	elapsedtime := 0
	for {
		tpsStart = transcount
		time.Sleep(time.Second)
		tpsEnd = transcount

		fmt.Printf("%d Running at %d TPS, %d trans\n", elapsedtime, tpsEnd-tpsStart, transcount)
		elapsedtime++
	}
}

func printEndpointStats() {
	for endpoint := range endpointTimes {
		responseTimes := endpointTimes[endpoint]
		totalTime, _ := time.ParseDuration("0")
		numCommands := float64(len(responseTimes))

		for _, timeDur := range responseTimes {
			totalTime += timeDur
		}

		// Some gross time type conversions
		avgTimeFloat := totalTime.Seconds() / numCommands
		avgTimeString := strconv.FormatFloat(avgTimeFloat, 'f', 6, 64)
		avgTimeTime, _ := time.ParseDuration(avgTimeString + "s")

		fmt.Printf("%v: %v\n", endpoint, avgTimeTime)
	}
}

func saveEndpointStats() {
	f, err := os.Create("./endpointStats.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for endpoint := range endpointTimes {
		responseTimes := endpointTimes[endpoint]
	}
}
