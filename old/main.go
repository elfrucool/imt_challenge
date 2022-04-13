package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Options = struct {
	SrcUrl       string
	DestFileName string
	BandWidthKbS uint
}

func main() {
	options, err := parseArguments()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: got error: %s\n", err.Error())
		usage()
		os.Exit(1)
	}

	fmt.Printf("Options: %v\n", options)

	reader, err := startDownload(options.SrcUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error when downloading url: %s, error: %s\n", options.SrcUrl, err.Error())
		os.Exit(2)
	}
	defer reader.Close()

	file, err := os.Create(options.DestFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error when creating file: %s, error: %s\n", options.DestFileName, err.Error())
		os.Exit(2)
	}
	defer file.Close()

	hash, err := computeHash(reader, options.BandWidthKbS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error while processing file contents, error: %s\n", err.Error())
		os.Exit(2)
	}

	hashStr := hex.EncodeToString(hash)

	fmt.Printf("Got hash value for %v => %s\n", hash, hashStr)

	file.WriteString(hashStr)
}

func parseArguments() (Options, error) {
	realArgs := os.Args[1:]
	if len(realArgs) < 2 {
		return Options{}, errors.New("Not enough arguments")
	}

	bWidthAsString := "0"
	if len(realArgs) > 2 {
		bWidthAsString = realArgs[2]
	}

	bWidthAsUint, err := strconv.ParseUint(bWidthAsString, 10, 16)
	if err != nil {
		return Options{}, errors.New(
			fmt.Sprintf("Invalid bandwidth, must be a number: %s", err.Error()),
		)
	}

	srcUrl := realArgs[0]
	destFileName := realArgs[1]

	return Options{
		SrcUrl:       srcUrl,
		DestFileName: destFileName,
		BandWidthKbS: uint(bWidthAsUint),
	}, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: %s <source-url> <dest-file> [<bandwidth-in-kb/s>]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "ARGUMENTS:")
	fmt.Fprintln(os.Stderr, "    <source-url>:")
	fmt.Fprintln(os.Stderr, "        The url to download the file")
	fmt.Fprintln(os.Stderr, "    <dest-file>:")
	fmt.Fprintln(os.Stderr, "        The destination file to store the hash (will be overwritten)")
	fmt.Fprintln(os.Stderr, "    <bandwidth-in-kb/s> (optional, default to 0):")
	fmt.Fprintln(os.Stderr, "        The bandwidth in kb/s for throttling (0 means no throttling)")
}

func startDownload(srcUrl string) (io.ReadCloser, error) {
	client := http.Client{}

	response, err := client.Get(srcUrl)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, errors.New(
			fmt.Sprintf("Got invalid status code: %d", response.StatusCode),
		)
	}

	return response.Body, nil
}

func computeHash(reader io.Reader, bandwidth uint) ([]byte, error) {
	var kbsRead uint = 0
	buffer := make([]byte, 1024) // 1kb block
	hash := make([]byte, 8)

	for {
		bytesRead, err := reader.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		for _, b := range buffer[:bytesRead] {
			imtHash(hash, b)
		}

		kbsRead++
		if bandwidth > 0 && (kbsRead%bandwidth == 0) {
			time.Sleep(1 * time.Second)
		}

		fmt.Printf("Downloaded %d kb\r", kbsRead)
	}

	fmt.Println()

	return hash, nil
}

func imtHash(hash []byte, b byte) {
	coefficients := []int{2, 3, 5, 7, 11, 13, 17, 19}

	for i := range hash {
		prev := 0
		if i > 0 {
			prev = int(hash[i-1])
		}

		current := int(b)

		newVal := ((prev + current) * coefficients[i]) % 255

		hash[i] = byte(newVal)
	}
}
