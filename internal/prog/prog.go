package prog

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elfrucool/imt_challenge/internal/args"
	"github.com/elfrucool/imt_challenge/internal/net"
	"github.com/elfrucool/imt_challenge/internal/process"
)

type WaitOneSecondThrottle struct{}

func (_ *WaitOneSecondThrottle) Throttle() {
	time.Sleep(1 * time.Second)
}

func Run(options args.Options) error {
	reader, err := net.StartDownload(options.SourceUrl)
	if err != nil {
		return err
	}
	defer reader.Close()

	file, err := os.Create(options.DestinationFile)
	if err != nil {
		return err
	}
	defer file.Close()

	hash, err := process.Process(
		file, &WaitOneSecondThrottle{}, options.BandWithInBytesPerSec(),
	)

	hashHexString := hex.EncodeToString(hash)

	log.Printf("[Run] got hash: %v (%v)", hashHexString, hash)

	file.WriteString(hashHexString)

	return nil
}

func Usage() {
	fmt.Fprintf(os.Stderr, "USAGE: %s <source-url> <dest-file> [<bandwidth-in-kb/s>]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "ARGUMENTS:")
	fmt.Fprintln(os.Stderr, "    <source-url>:")
	fmt.Fprintln(os.Stderr, "        The url to download the file")
	fmt.Fprintln(os.Stderr, "    <dest-file>:")
	fmt.Fprintln(os.Stderr, "        The destination file to store the hash (will be overwritten)")
	fmt.Fprintln(os.Stderr, "    <bandwidth-in-kb/s> (optional, default to 0):")
	fmt.Fprintln(os.Stderr, "        The bandwidth in kb/s for throttling (0 means no throttling)")
}
