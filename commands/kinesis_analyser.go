package commands

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/kj187/aws-utility/aws/kinesis"
	"github.com/fatih/color"
	"github.com/apcera/termtables"
	"github.com/gosuri/uiprogress"
)

func init() {
	RootCmd.AddCommand(kinesisStreamAnalyserCmd)
}

var kinesisStreamAnalyserCmd = &cobra.Command{
	Use: "kinesis:stream:analyse",
	Short: "Analyse a Kinesis stream",
	Run: func(cmd *cobra.Command, args []string) {
		streamAnalyse()
	},
}

func streamAnalyse() {
	fmt.Println("AWS Kinesis Stream Analyser\n")
	fmt.Println("-----------------------------------------------------\n")

	fmt.Println("Available streams:")
	streams, keys := kinesis.GetStreams()
	for _, k := range keys {
		key := color.GreenString("[" + strconv.Itoa(k) + "]")
		fmt.Println(key, streams[k])
	}

	var input int
	fmt.Print("\nPlease select a stream: ")
	fmt.Scanln(&input)

	if _, ok := streams[input]; ok { } else {
		fmt.Println("Stream not available!")
		return
	}

	fmt.Println("\nSelected stream: ", streams[input])

	shardIds, keys := kinesis.GetShardIds(streams[input])

	fmt.Println("Shards in stream: ", len(shardIds))
	fmt.Println("")

	table := termtables.CreateTable()
	table.AddHeaders("ShardId", "Records")

	// Progressbar
	bar := uiprogress.AddBar(len(shardIds)).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		current := b.Current()
		if (current == 0) { current = 1 }
		return fmt.Sprintf("Reading shard (%d/%d)", current, len(shardIds))
	})
	uiprogress.Start()
	var wg sync.WaitGroup

	var recordsInStream int = 0
	for _, k := range keys {
		wg.Add(1)

		recordCount := kinesis.GetShardRecordCount(shardIds[k], streams[input])
		recordsInStream = recordsInStream + recordCount
		table.AddRow(shardIds[k], recordCount)

		go func() {
			defer wg.Done()
			bar.Incr()
		}()
	}

	time.Sleep(time.Second) // wait for a second for all the go routines to finish
	wg.Wait()
	uiprogress.Stop()

	table.AddRow("", "---------")
	table.AddRow("", recordsInStream)

	fmt.Println("")
	fmt.Println(table.Render())
}