package Node

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/evscott/Distributed-NACN/Models"
	"github.com/evscott/Distributed-NACN/constants"
)

type Info struct {
	IP            string         `json:"IP"`
	Port          string         `json:"port"`
	Nodes         []string       `json:"nodes"`
	RequestBy     map[string]int `json:"requestBy"`
	Interested    bool           `json:"interested"`
	ObjectPresent bool           `json:"objectPresent"`
	Object        *Models.Object   `json:"Object"`
}

// Create is used a constructor that instantiates a new node using it's initial knowledge.
//
// A node must be created with initial knowledge of it's network IP, ID, and the IDs of it's neighbors.
func Create(ip, port string, nodes []string) *Info {
	rand.Seed(time.Now().UTC().UnixNano())
	newNode := Info{
		IP:   ip,
		Port: port,
		Nodes: nodes,
	}
	newNode.RequestBy = make(map[string]int)

	return &newNode
}

// Create is used a constructor that instantiates a new node using it's initial knowledge.
//
// A node must be created with initial knowledge of it's network IP, ID, and the IDs of it's neighbors.
func CreateWithObject(ip, port string, nodes []string) *Info {
	rand.Seed(time.Now().UTC().UnixNano())
	// Initialize object with empty map
	object := &Models.Object{}
	object.Obtained = make(map[string]int)

	newNode := Info{
		IP:   ip,
		Port: port,
		Nodes: nodes,
		Object: object,
		ObjectPresent: true,
	}
	newNode.RequestBy = make(map[string]int)

	return &newNode
}


func (i *Info) AcquireObject() {

}

func (i *Info) ReleaseObject() {

}

func (i *Info) request(reqSource string) {

}

func (i *Info) receiveObject(object *Models.Object) {
	i.Object = object
	i.ObjectPresent = true

	// Simulation an amount of time that a task involving the object might take
	time.Sleep(time.Millisecond)
}

// SendMsg handles sending messages across the distributed system using a destination.
func (i *Info) sendMsg(msg Models.Message, dest string) error {
	connOut, err := net.DialTimeout("tcp", i.IP+":"+dest, time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Couldn't send go to %s:%s \n", i.IP, dest)
			return err
		}
	}

	if err := json.NewEncoder(connOut).Encode(&msg); err != nil {
		fmt.Printf("Couldn't enncode message %v \n", msg)
		return err
	}
	return nil
}

// ListenOnPort is the communication satellite for a node that listens for incoming messages.
// Incoming messages are marshalled into a `Message` struct, and are directed to a handler
// depending on the messages `Intent`.
// Incoming messages that cannot be marshalled into a `Message` may cause erroneous behaviour.
func (i *Info) ListenOnPort() {
	ln, err := net.Listen("tcp", fmt.Sprint(":"+i.Port))
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("Starting node on %s:%s...\n", i.IP, i.Port)

	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				fmt.Printf("Error received while listening %s:%s \n", i.IP, i.Port)
			}
		}

		var msg Models.Message
		if err := json.NewDecoder(connIn).Decode(&msg); err != nil {
			fmt.Printf("Error decoding %v\n", err)
		}

		switch msg.Intent {
		case constants.IntentSendObject:
			i.receiveObject(msg.Object)
		case constants.IntentRequestObject:
			i.request(msg.Source)
		}
	}
}
