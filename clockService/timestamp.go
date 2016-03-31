package clockService

// Timestamp is one of most important feature in distributed system.
// It helps remote nodes to know what event happens first, or how should they
// order the events.
// The usage of timestamp is highly related with the Clock Service built
// on top of it, for more information, see clockService.go
type Timestamp interface {
	// Compare with another timestamp.
	//
	// @param: input: another timestamp
	// @return: an integer indicating which one is smaller, notice
	//          this doesn't guarantee one event happens before another,
	//          that depends on the type of clock service built on top of
	//          the timestamp
	CompareTo(Timestamp) int

	// Convert timestamp to string, so could be packed in a message.
	String() string
}
