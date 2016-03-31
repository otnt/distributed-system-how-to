package clockService

import (
	"fmt"
)

// Create a new ClockService type
// currently, only supporting logical clock
//
// @param: clockType string
//         any of following(seperated by comma): logical
//         this function will panic if input is any other value
// @return: ClockService as required by user
func NewClockService(clockType string) ClockService {
	switch clockType {
	case "logical":
		return &LogicalClock{0}
	default:
		panic(fmt.Sprintf("Clock Service type %s is not supported", clockType))
	}
}

// A ClockService serves as part of main infrastructure, helping all components
// to know the happen sequence of any two events. It trys to provide a mechanism
// that could help determine whether two events are causally related (in a best effort).
// It also takes important role to provide a good order of events(FIFO, Causal, Total etc.),
// but the ordering is mainly provided by other service, such as an ordering service.
//
// Several clock is widely used in distributed systems. Lamport Clock/Logical Clock is
// no doubt the simplest one, but could handle a majority of situations. Vector Clock could
// provide more information to help determine the happen sequence of events, but with tradeoff
// of more network overhead. Bounded Vector Clock is a tradeoff between Logical Clock and Vector
// Clock, i.e. it balance the network overhead and ability to determince event happen sequence.
// Dependency Clock is more complicated, it provides same ability of sequence determine as
// Vector Clock, but also has constant network overhead, with tradeoff of more sophasticated
// implementation.
//
// More info:
// Fundamentals of Distributed Computing: A Practical Tour of Vector Clock Systems
type ClockService interface {
	// Get a new Timestamp from ClockService. This function returns current/most up-to-date
	// time, and then increment current time.
	NewTimestamp() Timestamp

	// Update time given a receiving Timestamp. This function is used when receiving a new
	// message, so we need to update our local time with tiem on other nodes.
	UpdateTimestamp(Timestamp)
}
