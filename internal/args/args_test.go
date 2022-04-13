package args

import (
	"log"
	"testing"
)

// Tests for missing arguments
func TestParseArguments_NotEnoughArguments(t *testing.T) {
	t.Run("no-arguments", func(t *testing.T) {
		_, err := ParseArguments([]string{"the-program"})
		if assertError(t, err) {
			return
		}
	})

	t.Run("only-one-argument", func(t *testing.T) {
		_, err := ParseArguments([]string{"the-program", "the-url"})
		if assertError(t, err) {
			return
		}
	})
}

// Tests for invalid bandwidth with some sample data.
func FuzzInvalidBandWidth(f *testing.F) {
	for _, testCase := range []string{"foo", "bar", "baz", "-1"} {
		f.Add(testCase)
	}

	f.Fuzz(func(t *testing.T, badBandWidth string) {
		log.Printf("[FussInvalidBandWidth] badBandWidth: %s", badBandWidth)
		_, err := ParseArguments([]string{"the-program", "the-url", "the-filename", badBandWidth})
		if assertError(t, err) {
			return
		}
	})
}

// Test for valid arguments: without bandwidth
func TestValidArgumentsNoBandWidth(t *testing.T) {
	options, err := ParseArguments([]string{"the-program", "the-url", "the-filename"})

	if assertNoError(t, err) {
		return
	}

	if assertEquals(t, options.SourceUrl, "the-url") {
		return
	}

	if assertEquals(t, options.DestinationFile, "the-filename") {
		return
	}

	if assertEquals(t, options.BandWidthKbS, 0) {
		return
	}
}

// Test for valid arguments: with bandwidth
func TestValidArgumentsBandWidth(t *testing.T) {
	options, err := ParseArguments([]string{"the-program", "the-url", "the-filename", "10"})

	if assertNoError(t, err) {
		return
	}

	if assertEquals(t, options.SourceUrl, "the-url") {
		return
	}

	if assertEquals(t, options.DestinationFile, "the-filename") {
		return
	}

	if assertEquals(t, options.BandWidthKbS, 10) {
		return
	}
}

// Passes if err is not nil.
// Returns true if there was an unexpected assertion to early return.
func assertError(t *testing.T, err error) bool {
	if err == nil {
		t.Error("Expecting to have an error")
		return true
	}
	return false
}

// Passess if err is nil
// Returns true if there was an unexpected assertion to early return.
func assertNoError(t *testing.T, err error) bool {
	if err != nil {
		t.Errorf("Expecting not to have an error, got %v", err)
		return true
	}
	return false
}

// Asserts if a value is equal to other value
// Returns true if there was an unexpected assertion to early return.
func assertEquals[T comparable](t *testing.T, actual T, expected T) bool {
	if actual != expected {
		t.Errorf("Expeting %v to be %v", actual, expected)
		return true
	}
	return false
}
