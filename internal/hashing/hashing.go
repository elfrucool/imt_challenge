package hashing

import (
	"errors"
)

const (
	HashBufferSize int = 8
)

var (
	coefficients = []int{2, 3, 5, 7, 11, 13, 17, 19}
)

func Coefficients() []int {
	c := make([]int, len(coefficients))
	copy(c, coefficients)
	return c
}

// the hasing function for each byte
// it modifies in place the provided hash
// this is to save space
//
// *Note:* we are assuming the donwloaded file may not fit
// into memory nor hard disk
// It seems to be more efficient to modify in place
// instead of creating new hashes and allowing
// garbage collector to clean up unused variables.
func ImtHash(hash []byte, b byte) error {
	if len(hash) != HashBufferSize {
		return errors.New("Invalid hash buffer size, it must be 8")
	}

	// based on tests I discovered previous hash values do not matter at all
	var prev int
	for i := 0; i < len(hash); i++ {
		current := int(b)

		newVal := ((int(prev) + current) * coefficients[i]) % 255

		hash[i] = byte(newVal)

		prev = newVal
	}

	return nil
}
