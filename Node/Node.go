package Node

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/evscott/Distributed-ANA/Models"
	"github.com/evscott/Distributed-ANA/constants"
)

type Info struct {
	IP            string         `json:"IP"`
	Port          string         `json:"port"`
	RequestBy     map[string]int `json:"requestBy"`
	Interested    bool           `json:"interested"`
	ObjectPresent bool           `json:"objectPresent"`
	Object        *Models.Object `json:"object"`
	Parent        string         `json:"parent"`
	Next          *string        `json:"next"`
}

// Create is used a constructor that instantiates a new node using it's initial knowledge.
//
// A node must be created with initial knowledge of it's network IP, ID, and the IDs of it's neighbors.
func Create(ip, port, parent string) *Info {
	rand.Seed(time.Now().UTC().UnixNano())
	newNode := Info{
		IP:     ip,
		Port:   port,
		Parent: parent,
	}
	newNode.RequestBy = make(map[string]int)

	return &newNode
}

// Create is used a constructor that instantiates a new node using it's initial knowledge.
//
// A node must be created with initial knowledge of it's network IP, ID, and the IDs of it's neighbors.
func CreateWithObject(ip, port, parent string) *Info {
	rand.Seed(time.Now().UTC().UnixNano())
	// Initialize object with empty map
	object := &Models.Object{}
	object.Obtained = make(map[string]int)

	newNode := Info{
		IP:            ip,
		Port:          port,
		Object:        object,
		ObjectPresent: true,
		Parent:        parent,
	}
	newNode.RequestBy = make(map[string]int)

	return &newNode
}

func (i *Info) AcquireObject() {
	i.Interested = true
	if !i.ObjectPresent {
		msg := Models.Message{
			Source: i.Port,
			Intent: constants.IntentRequestObject,
		}
		if err := i.sendMsg(msg, i.Parent); err != nil {
			fmt.Printf("Error requestion object: %v\n", err)
		}
		i.Parent = i.Port
	}
}

func (i *Info) ReleaseObject() {
	i.Interested = false
	if i.Next != nil {
		msg := Models.Message{
			Source: i.Port,
			Intent: constants.IntentSendObject,
			Object: i.Object,
		}
		if err := i.sendMsg(msg, *i.Next); err != nil {
			fmt.Printf("Error sending object after release: %v\n", err)
		}
		i.ObjectPresent = false
		i.Next = nil
	}
}

func (i *Info) request(reqSource string) {
	if i.Parent != i.Port {
		msg := Models.Message{
			Source: reqSource,
			Intent: constants.IntentRequestObject,
		}
		if err := i.sendMsg(msg, i.Parent); err != nil {
			fmt.Printf("Error passing on request: %v\n", err)
		}
	} else if i.Interested {
		msg := Models.Message{
			Source: i.Port,
			Intent: constants.IntentSendObject,
			Object: i.Object,
		}
		if err := i.sendMsg(msg, reqSource); err != nil {
			fmt.Printf("Error sending object to requester: %v\n", err)
		}
		i.ObjectPresent = false
	}
	i.Parent = reqSource
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
