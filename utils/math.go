package utils

// Return maximum 32-bit integer of two input 32-bit integer
func MaxInt32(x int, y int) int {
	if x >= y {
		return x
	} else {
		return y
	}
}

// Return maximum 64-bit integer of two input 64-bit integer
func MaxInt64(x int64, y int64) int64 {
	if x >= y {
		return x
	} else {
		return y
	}
}
