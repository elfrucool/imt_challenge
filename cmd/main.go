package main

import (
	"log"
	"os"

	"github.com/elfrucool/imt_challenge/internal/args"
	"github.com/elfrucool/imt_challenge/internal/prog"
)

func main() {
	options, err := args.ParseArguments(os.Args)

	if err != nil {
		log.Printf("[main] error parsing arguments: %v", err)
		prog.Usage()
		os.Exit(1)
	}

	log.Printf("[main] running program with options: %v", options)

	err = prog.Run(options)

	if err != nil {
		log.Fatalf("[main] error processing file: %v", err)
	}
}
