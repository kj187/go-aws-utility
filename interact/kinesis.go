package interact

import (
	"fmt"
	"strconv"
	"errors"

	"github.com/fatih/color"
	"github.com/kj187/aws-utility/aws/kinesis"
)

func AskForKinesisStream() (string, error) {
	fmt.Println("Available streams:")
	streams, keys := kinesis.GetStreams()

	if len(streams) <= 0 {
		return "", errors.New(color.YellowString("WARNING: No streams available!"))
	}

	for _, k := range keys {
		key := color.GreenString("[" + strconv.Itoa(k) + "]")
		fmt.Println(key, streams[k])
	}

	var input int
	fmt.Print("\nPlease select a stream: ")
	fmt.Scanln(&input)

	if _, ok := streams[input]; ok { } else {
		return "", errors.New(color.RedString("ERROR: Stream not available!"))
	}

	fmt.Println("\nSelected stream: ", streams[input])

	return streams[input], nil
}
