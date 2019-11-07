package main

import (
	"fmt"
	"github.com/evscott/Distributed-NACN/Node"
	"net"
	"strings"
	"time"
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
	nodes := []string{"8001", "8002", "8003"}

	// Create & start nodes with initial knowledge
	n1 := Node.Create(ip, "8001", nodes)
	n2 := Node.CreateWithObject(ip, "8002", nodes)
	n3 := Node.Create(ip, "8003", nodes)
	go n1.ListenOnPort()
	go n2.ListenOnPort()
	go n3.ListenOnPort()
	time.Sleep(time.Millisecond)

	go logSystemTraffic(n1, n2, n3)

	n1.AcquireObject()
	time.Sleep(time.Millisecond)
	n3.AcquireObject()
	time.Sleep(time.Millisecond)
	n1.ReleaseObject()
	time.Sleep(time.Second)
}
