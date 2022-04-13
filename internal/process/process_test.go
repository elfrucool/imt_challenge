package process

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/elfrucool/imt_challenge/internal/hashing"
)

var i io.Reader

type CountingThrottling struct {
	total uint
}

func (c *CountingThrottling) Throttle() {
	c.total++
	log.Printf("[CountingThrottling] throttled %d times", c.total)
}

// Empty reader should return blank hash
func TestProcessEmptyReader(t *testing.T) {
	var throttling *CountingThrottling = &CountingThrottling{}

	hash, err := Process(bytes.NewReader(make([]byte, 0)), throttling, 0)

	if assertNoError(t, err) {
		return
	}

	if assertHashEqualTo(t, hash, make([]byte, hashing.HashBufferSize)) {
		return
	}

	if assertNoThrottling(t, throttling) {
		return
	}
}

// single byte in reader means single call of hashing function
func TestProcessSingleByteReader(t *testing.T) {
	var throttling *CountingThrottling = &CountingThrottling{}

	hash, err := Process(bytes.NewReader([]byte{12}), throttling, 0)

	if assertNoError(t, err) {
		return
	}

	if assertHashEqualTo(t, hash, []byte{24, 108, 90, 204, 81, 189, 102, 126}) {
		return
	}

	if assertNoThrottling(t, throttling) {
		return
	}
}

// Should process bytes, but since last one is zero, then hash will be zero
func TestProcessMoreBytes_FinalOneIsZero(t *testing.T) {
	throtling := &CountingThrottling{}

	hash, err := Process(bytes.NewReader([]byte{4, 12, 0}), throtling, 0)

	if assertNoError(t, err) {
		return
	}

	if assertHashEqualTo(t, hash, make([]byte, hashing.HashBufferSize)) {
		return
	}

	if assertNoThrottling(t, throtling) {
		return
	}
}

func TestProcess_WithThrottling(t *testing.T) {
	throttling := &CountingThrottling{}

	hash, err := Process(bytes.NewReader([]byte{0, 0, 0}), throttling, 1)

	if assertNoError(t, err) {
		return
	}

	if assertHashEqualTo(t, hash, make([]byte, hashing.HashBufferSize)) {
		return
	}

	if assertThrottlingWasCalledTimes(t, 3, throttling) {
		return
	}
}

// asserts if actual is equal to an expected hash.
// Returns true in case of error for early returns
func assertHashEqualTo(t *testing.T, actual []byte, expected []byte) bool {
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expecting hash to be %v, it is %v", expected, actual)
		return true
	}
	return false
}

// asserts the throttling mechanism was not called
// Returns true in case of error for early returns
func assertNoThrottling(t *testing.T, throttling *CountingThrottling) bool {
	if throttling.total > 0 {
		t.Errorf("throttling should not be called, it was called %d times", throttling.total)
		return true
	}
	return false
}

// Asserts the throttling mechanism was called a certain number of times
// Returns true in case of error for early returns
func assertThrottlingWasCalledTimes(t *testing.T, expectedTimes uint, throttling *CountingThrottling) bool {
	if expectedTimes != throttling.total {
		t.Errorf("throttling should be called %d times, but was %d times", expectedTimes, throttling.total)
		return true
	}
	return false
}

// Asserts that there is no error.
// Returns true in case of error for early returns
func assertNoError(t *testing.T, err error) bool {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return true
	}
	return false
}
