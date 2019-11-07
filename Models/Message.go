package Models

import (
	"github.com/evscott/Distributed-NACN/constants"
)

// Pair representing a pair of node ID's, I and J.
type Object struct {
	Obtained map[string]int `json:"obtained"`
}

// The format for a Request/Response in finding a shortest path across a distributed system.
//
// Source represents the Message sender.
// Intent represents the messages intent; i.e., whether it is to be handled by `Update` or some other handler.
type Message struct {
	Source string           `json:"source"`
	Intent constants.Intent `json:"intent"`
	Object *Object           `json:"object"`
}

// Just for pretty printing Request/Response info.
func (req Message) String() string {
	return "Message:{ Origin:" + req.Source + ", Intent: " + string(req.Intent) + " }\n"
}
