package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

func print_open_port(host string, port int) {
	var hostWithPort string
	hostWithPort = net.JoinHostPort(host, strconv.Itoa(port))
	_, err := net.DialTimeout("tcp", hostWithPort, 1 * 1000 * 1000 * 1000) // timeout 1s
	if err == nil {
		fmt.Printf("port %d\n", port)
	}
}

func connect_port_interval(wg *sync.WaitGroup, host string, start int, end int) {
	defer wg.Done()
	var port int
	for port = start; port < end; port++ {
		print_open_port(host, port)	
	}
}

func main() {
	fmt.Println(len(os.Args))

	var wg sync.WaitGroup

	var host string = "127.0.0.1"

	var numThreads int = 20;
	var maxNumPorts int = 60000
	var portsPerThread int = maxNumPorts/numThreads

	for thread := 0; thread < numThreads; thread++ {
		wg.Add(1)
		var start int = thread * portsPerThread
		go connect_port_interval(&wg, host, start, start + portsPerThread)
	}
	wg.Wait()
}
