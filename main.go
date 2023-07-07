package main

import (
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

func get_ports_list_from_args() []int {
	var portsArray []int
	var portsStringArray []string = strings.Split(os.Args[2], ",")
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

func main() {
	var numArgs = len(os.Args)

	var host string

	//TODO 	verify num threads if num ports is less than num threads,
	// 		check if it is negative or 0
	//TODO 	check timeout, if server do not exist or is not reacheable
	// 		it takes too long
	var numThreads int = 1
	var maxNumPorts int = 65535 // 1 - 65535

	var lastArg int = 1

	var portsArray []int

	if numArgs > 1 {
		if os.Args[1] == "-p" {
			if numArgs != 4 {
				print_usage_and_exit()
			}
			lastArg = 3
			portsArray = get_ports_list_from_args()
			host = os.Args[lastArg]
		} else {
			host = os.Args[lastArg]
			portsArray = create_array_from_range(1, maxNumPorts + 1)
		}
	} else {
		print_usage_and_exit()
	}

	var wg sync.WaitGroup

	var numPorts int = len(portsArray)
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
