package args

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/elfrucool/imt_challenge/internal/common"
)

type Options struct {
	SourceUrl       string
	DestinationFile string
	BandWidthKbS    common.KbsPerSecond
}

func (o Options) BandWithInBytesPerSec() common.BytesPerSecond {
	return common.BytesPerSecond(o.BandWidthKbS * 1024)
}

func ParseArguments(args []string) (Options, error) {
	realArgs := args[1:]
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
		SourceUrl:       srcUrl,
		DestinationFile: destFileName,
		BandWidthKbS:    common.KbsPerSecond(bWidthAsUint),
	}, nil
}
