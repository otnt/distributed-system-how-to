package clockService

import (
	utils "github.com/otnt/distributed-system-notes/utils"
)

// Logical Clock is also called Lamport Clock. The concept of this type of clock is quite
// simple, it use a single integer as clock. So that the larger the integer is, the more
// possible that this event happens later. Notice Logic Clock does not guarantee the larger
// time means causally happen later.
//
// How this clock works?
// Initially, the time is 0. There are three actions defined with logical clock:
// Let's assumet the local time is T, so initially T = 0
// 1.When an internal event happens, T => T+1
// 2.When send a message m out, T is attached to m, then T => T+1
// 3.When receive a message m, let T' denotes time attached with m, then T => max(T, T') + 1
//
// What can Logical Clock guarantee?
// It guarantees that if two events e1, e2 with time T1, T2 respectively, then
// e1 -> e2 => T1 < T2, where -> is the 'happen before' operation defined by Lamport
// The equivalence of the former formula is,
// T1 >= T2 => e1 -> e2 is impossible
// This means if event e1 has time larger than or equal to event e2, then e1 could not
// happen before e2.
//
// What Logical Clock could not provide?
// Actually, logical clock could not determine any two event that which one happens before
// the other, or they are in parallel.
//
// Is Logical Clock enough for my system?
// It depends...
// The limitation of Logical Clock mostly lays in that it could not determine the happen
// sequence of any two events. This does not effect in most of systems.
// IMHO, for a system where no network partition may happen, or network partition would not make
// any bad effect, then generally it is enough to use logical clock. However, if network
// partition may result in data inconsistency, and the system you are implementing need
// to make sure data would be consistent eventually or all the time, then logical clock
// may not fit you need. In that case, you may find Vector Clock or Dependency Clock more
// pratical.
type LogicalClock struct {
	currentTimestamp int64
}

// Create a new logical clock.
func NewLogicalClock() LogicalClock {
	return LogicalClock{0}
}

// Get a new timestamp from logical clock.
func (lc *LogicalClock) NewTimestamp() (ts Timestamp) {
	ts = &LogicalTimestamp{lc.currentTimestamp}
	lc.currentTimestamp++
	return
}

// Update local time using timestamp in receiving message from other nodes.
func (lc *LogicalClock) UpdateTimestamp(timestamp Timestamp) {
	lt := timestamp.(*LogicalTimestamp)
	lc.currentTimestamp = utils.MaxInt64(lc.currentTimestamp, lt.time)
	lc.currentTimestamp++
}
