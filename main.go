package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func print_open_port(channel chan int, host string, port int) {
	var hostWithPort string
	hostWithPort = net.JoinHostPort(host, strconv.Itoa(port))
	_, err := net.DialTimeout("tcp", hostWithPort, 1 * 1000 * 1000 * 1000) // timeout 1s
	if err == nil {
		channel <- port
	}
}

func connect_port_interval(wg *sync.WaitGroup, channel chan int, host string, portsArray []int, start int, end int) {
	defer wg.Done()
	var port int
	for port = start; port < end; port++ {
		print_open_port(channel, host, portsArray[port])	
	}
}

func create_array_from_range(start int, end int) []int {
	var array []int = make([]int, end-start-1)
	for i := range array {
		array[i] = start + i
	}
	return array
}

func get_ports_list_from_args(pos int) []int {
	var portsArray []int
	var portsStringArray []string = strings.Split(os.Args[pos], ",")
	portsArray = make([]int, len(portsStringArray))
	for i, p := range portsStringArray {
		portsArray[i], _ = strconv.Atoi(p)
	}
	return portsArray
}

func print_usage_and_exit() {
	fmt.Println("Usage: program [options] host")
	os.Exit(1)
}

func getArgPositionIfExist(args []string, option string) (int, error) {
	for i, val := range args {
		if val == option {
			return i, nil
		}
	}
	return -1, errors.New("option not found")
}

func verifyNumberOfArgs(args []string) {
	var numArgs = len(args)
	var expectedNumArgs int = 2
	var err error
	_, err = getArgPositionIfExist(os.Args, "-t")
	if err == nil {
		expectedNumArgs += 2
	}
	_, err = getArgPositionIfExist(os.Args, "-p")
	if err == nil {
		expectedNumArgs += 2
	}
	if numArgs != expectedNumArgs {
		print_usage_and_exit()
	}
}

func main() {
	var host string

	//TODO 	verify num threads if num ports is less than num threads,
	// 		check if it is negative or 0
	//TODO 	check timeout, if server do not exist or is not reacheable
	// 		it takes too long
	var numThreads int = 1
	var maxNumPorts int = 65535 // 1 - 65535

	var lastArg int = 1

	var portsArray []int

	verifyNumberOfArgs(os.Args) // exit program if wrong number of args

	threadArgPosition, err := getArgPositionIfExist(os.Args, "-t")
	if err == nil {
		lastArg += 2
		numThreads, _ = strconv.Atoi(os.Args[threadArgPosition+1])
		if numThreads <= 0 {
			fmt.Println("Number of threads less than or equal to 0")
			os.Exit(2)
		}
	}
	portsArgPosition, err := getArgPositionIfExist(os.Args, "-p")
	if err == nil {
		lastArg += 2
		portsArray = get_ports_list_from_args(portsArgPosition+1)
	} else {
		portsArray = create_array_from_range(1, maxNumPorts + 1)
	}

	host = os.Args[lastArg]

	var wg sync.WaitGroup

	var numPorts int = len(portsArray)
	if numPorts < numThreads {
		numThreads = numPorts
	}
	var portsPerThread int = numPorts/numThreads

	channel := make(chan int, numPorts)

	fmt.Printf("Scanning ports from host %s...\n", host)

	for thread := 0; thread < numThreads; thread++ {
		wg.Add(1)
		var start int = thread * portsPerThread
		go connect_port_interval(&wg, channel, host, portsArray, start, start + portsPerThread)
	}
	wg.Wait()

	var outputPorts []int

	close(channel)
	for elem := range channel {
		outputPorts = append(outputPorts, elem)
	}

	sort.Ints(outputPorts)
	for _, elem := range outputPorts {
		fmt.Printf("port %d\n", elem)
	}
}
