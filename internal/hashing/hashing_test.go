package hashing

import (
	"bytes"
	"fmt"
	"testing"
)

func TestImtHashInvalidSize(t *testing.T) {
	hash := make([]byte, 0)

	err := ImtHash(hash, 0)

	if err == nil {
		t.Error("Expecting function to return an error")
	}
}

// Testing the example provided in challenge
func TestBasicExcample(t *testing.T) {
	hash := make([]byte, HashBufferSize)
	expected := []byte{24, 108, 90, 204, 81, 189, 102, 126}

	ImtHash(hash, 12)

	if !bytes.Equal(hash, expected) {
		t.Errorf("Invalid hash, expected:\n %v\ngot:\n %v", expected, hash)
	}
}

// Following tests use an approach named property-based testing,
// see https://forallsecure.com/blog/what-is-property-based-testing
//
// A property is a predictable behavior that arises from the target logic
// based on changed a fixed set of inputs
// (e.g. modifying a single parameter while leaving the rest of parameters
// constant)
//
// TODO: use Fuzz tests: https://pkg.go.dev/testing#F.Fuzz

// For this case, we iterate over the "current byte"
// leaving the original hash with zeroes and assess the calculations
// behavior based only on the input byte.
func TestImtHashZeroesByteMultipliesByCoefficient(t *testing.T) {
	const maxByte int = 255

	coefficients := Coefficients()

	// iterating using int to avoid passing from 255 to 0 inside the loop
	for byteAsInt := 0; byteAsInt <= maxByte; byteAsInt++ {
		// using t.Run to label each individual scenario
		t.Run(fmt.Sprintf("byte=%d", byteAsInt), func(t2 *testing.T) {
			hash := make([]byte, HashBufferSize)

			currentByte := byte(byteAsInt)

			err := ImtHash(hash, currentByte)

			if err != nil {
				t2.Errorf("expected err to be null, got %v", err)
				return
			}

			// validate every value in the hash is as expected
			prevValue := 0
			for index, coefficient := range coefficients {
				expectedValue := coefficient * (int(byteAsInt) + prevValue) % 255

				if byte(expectedValue) != hash[index] {
					t2.Errorf(
						"[b=%d] Invalid byte at position %d, expecting %d, got %d",
						currentByte, index, expectedValue, hash[index])
					return
				}

				prevValue = expectedValue
			}
		})
	}
}

// We are now evaluating the function changing the values inside the previous
// hash buffer, leaving current byte as zero.
//
// We are going to set all bytes with the same value and evaluate how
// the value is used for next position in combination with coefficients.
func TestOnZeroByteEverythingEndsAzZero(t *testing.T) {
	const maxByte int = 255

	// iterating using int to avoid passing from 255 to 0 inside the loop
	for byteAsInt := 0; byteAsInt <= maxByte; byteAsInt++ {
		// using t.Run to label each individual scenario
		t.Run(fmt.Sprintf("byte=%d", byteAsInt), func(t2 *testing.T) {
			hash := bytes.Repeat([]byte{byte(byteAsInt)}, HashBufferSize)

			err := ImtHash(hash, 0)

			if err != nil {
				t2.Errorf("expected err to be null, got %v", err)
				return
			}

			if !bytes.Equal(hash, bytes.Repeat([]byte{0}, HashBufferSize)) {
				t2.Errorf("Expecting hash to be all zeroed, got: %v", hash)
				return
			}
		})
	}
}
