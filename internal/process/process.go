package process

import (
	"fmt"
	"io"
	"log"

	"github.com/elfrucool/imt_challenge/internal/common"
	"github.com/elfrucool/imt_challenge/internal/hashing"
)

// mechanism to slow down things in the process
type Throttling interface {
	Throttle()
}

// Processes the input in chunks until the end
// it calculates the hash and returns it
func Process(reader io.Reader, throttling Throttling, bandWidthInBytes common.BytesPerSecond) ([]byte, error) {
	bufferSize := 1024
	if bandWidthInBytes > 0 {
		bufferSize = int(bandWidthInBytes)
	}

	hash := make([]byte, hashing.HashBufferSize)
	buffer := make([]byte, bufferSize)
	totalBytesProcessed := uint64(0)

	for {
		bytesRead, err := reader.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		for _, currentByte := range buffer[:bytesRead] {
			hashing.ImtHash(hash, currentByte)
		}

		totalBytesProcessed += uint64(bytesRead)
		if bandWidthInBytes > 0 {
			throttling.Throttle()
		}

		// sending debug data to console to a single line
		fmt.Printf("[Process] Processed %d bytes             \r", totalBytesProcessed)
	}

	// moving to next line since \r (above) positions cursor in begining of an laready written line
	fmt.Println()

	log.Printf("[Process] Finish process, %d bytes read", totalBytesProcessed)
	return hash, nil
}
