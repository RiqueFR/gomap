package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sort"
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

func main() {
	fmt.Println(len(os.Args))

	var wg sync.WaitGroup

	var host string = "127.0.0.1"

	var numThreads int = 20
	var maxNumPorts int = 65535 // 1 - 65535

	var portsArray []int = create_array_from_range(1, maxNumPorts + 1)
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
