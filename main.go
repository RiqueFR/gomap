package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {
	var port int
	var host string = "127.0.0.1"

	var hostWithPort string

	for port = 1; port < 60000; port++ {
		hostWithPort = net.JoinHostPort(host, strconv.Itoa(port))
		_, err := net.DialTimeout("tcp", hostWithPort, 1 * 1000 * 1000 * 1000) // timeout 1s
		if err == nil {
			fmt.Printf("port %d\n", port)
		}
		
	}
}
