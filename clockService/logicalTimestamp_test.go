package clockService

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	_ = LogicalTimestamp{9}
}

func TestCompareTo(t *testing.T) {
	t1 := &LogicalTimestamp{0}
	t2 := &LogicalTimestamp{1}
	t3 := &LogicalTimestamp{2}

	//same value gave 0
	assert.Equal(t, t1.CompareTo(t1), 0)
	assert.Equal(t, t2.CompareTo(t2), 0)
	assert.Equal(t, t3.CompareTo(t3), 0)

	//small is -1
	assert.Equal(t, t1.CompareTo(t2), -1)
	assert.Equal(t, t2.CompareTo(t3), -1)

	//large is 1
	assert.Equal(t, t2.CompareTo(t1), 1)
	assert.Equal(t, t3.CompareTo(t2), 1)
}

func TestString(t *testing.T) {
	t1 := &LogicalTimestamp{0}
	assert.Equal(t, t1.String(), "0")
}
