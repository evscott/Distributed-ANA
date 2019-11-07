package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/evscott/Distributed-ANA/Node"
)

// main is the entry point for this distributed system.
func main() {
	ipAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Print(err)
		return
	}
	ip := strings.Split(ipAddr[0].String(), "/")[0]

	runExample(ip)
}

func logSystemTraffic(n1, n2, n3 *Node.Info) {
	if n1.ObjectPresent {
		fmt.Printf("Object at Node %s\n", n1.Port)
	} else if n2.ObjectPresent {
		fmt.Printf("Object at Node %s\n", n2.Port)
	} else if n3.ObjectPresent {
		fmt.Printf("Object at Node %s\n", n3.Port)
	} else {
		fmt.Printf("Object being processed... \n")
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
		if n1.ObjectPresent {
			fmt.Printf("Object at Node %s\n", n1.Port)
		} else if n2.ObjectPresent {
			fmt.Printf("Object at Node %s\n", n2.Port)
		} else if n3.ObjectPresent {
			fmt.Printf("Object at Node %s\n", n3.Port)
		}
	}
}

func runExample(ip string) {
	// Create & start nodes with initial knowledge
	n1 := Node.CreateWithObject(ip, "8001", "8001")
	n2 := Node.Create(ip, "8002", "8001")
	n3 := Node.Create(ip, "8003", "8001")
	n4 := Node.Create(ip, "8004", "8001")
	n5 := Node.Create(ip, "8005", "8001")
	go n1.ListenOnPort()
	go n2.ListenOnPort()
	go n3.ListenOnPort()
	go n4.ListenOnPort()
	go n5.ListenOnPort()
	time.Sleep(time.Millisecond)

	go logSystemTraffic(n1, n2, n3)

	time.Sleep(time.Second)
}
