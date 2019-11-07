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

func logSystemTraffic(n1, n2, n3, n4, n5 *Node.Info) {
	if n1.ObjectPresent {
		fmt.Printf("Object at Node %s\n", n1.Port)
	} else if n2.ObjectPresent {
		fmt.Printf("Object at Node %s\n", n2.Port)
	} else if n3.ObjectPresent {
		fmt.Printf("Object at Node %s\n", n3.Port)
	} else if n4.ObjectPresent {
		fmt.Printf("Object at Node %s\n", n2.Port)
	} else if n5.ObjectPresent {
		fmt.Printf("Object at Node %s\n", n3.Port)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
		if n1.ObjectPresent {
			fmt.Printf("Object at Node %s\n", n1.Port)
		} else if n2.ObjectPresent {
			fmt.Printf("Object at Node %s\n", n2.Port)
		} else if n3.ObjectPresent {
			fmt.Printf("Object at Node %s\n", n3.Port)
		} else if n4.ObjectPresent {
			fmt.Printf("Object at Node %s\n", n2.Port)
		} else if n5.ObjectPresent {
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

	go logSystemTraffic(n1, n2, n3, n4, n5)

	go n2.AcquireObject()
	time.Sleep(time.Millisecond)
	go n3.AcquireObject()
	time.Sleep(time.Millisecond)
	go n4.AcquireObject()
	time.Sleep(time.Millisecond)
	go n2.ReleaseObject()
	time.Sleep(time.Millisecond)
	go n2.AcquireObject()

	time.Sleep(time.Second)
	fmt.Printf("{ Node %s, Parent: %s, Next: %s }\n", n1.Port, n1.Parent, n1.Next)
	fmt.Printf("{ Node %s, Parent: %s, Next: %s }\n", n2.Port, n2.Parent, n2.Next)
	fmt.Printf("{ Node %s, Parent: %s, Next: %s }\n", n3.Port, n3.Parent, n3.Next)
	fmt.Printf("{ Node %s, Parent: %s, Next: %s }\n", n4.Port, n4.Parent, n4.Next)
	fmt.Printf("{ Node %s, Parent: %s, Next: %s }\n", n5.Port, n5.Parent, n5.Next)
}
