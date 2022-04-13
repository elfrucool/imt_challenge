package process

import (
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

	log.Printf(
		"[Process] reader: %v, throttling: %v, bandW: %d, bufferSize: %d",
		reader, throttling, bandWidthInBytes, bufferSize,
	)

	hash := make([]byte, hashing.HashBufferSize)
	buffer := make([]byte, bufferSize)
	totalBytesProcessed := uint64(0)

	for {
		bytesRead, err := reader.Read(buffer)

		log.Printf(
			"[Process] err: %v, bytesRead: %d, hash: %v, buffer: (%d bytes) %v",
			err, bytesRead, hash, len(buffer), buffer,
		)

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

		log.Printf("[Process] Processing, %d bytes", totalBytesProcessed)
	}

	log.Printf("[Process] Finish process, %d bytes read", totalBytesProcessed)
	return hash, nil
}
