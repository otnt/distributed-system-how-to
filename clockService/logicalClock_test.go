package clockService

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateLogicalClock(t *testing.T) {
	_ = NewLogicalClock()
}

func TestLogicalClockNewTimestamp(t *testing.T) {
	lc := NewLogicalClock()
	assert.Equal(t, lc.NewTimestamp(), &LogicalTimestamp{0})
}

func TestLogicalCLockUpdateTimestamp(t *testing.T) {
	lc := NewLogicalClock()
	lc.UpdateTimestamp(&LogicalTimestamp{10})
	assert.Equal(t, lc.NewTimestamp(), &LogicalTimestamp{11})
}
