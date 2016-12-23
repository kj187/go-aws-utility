package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/kj187/aws-utility/aws/kinesis"
	"github.com/apcera/termtables"
	"github.com/gosuri/uiprogress"
	"github.com/kj187/aws-utility/interact"
)

func init() {
	RootCmd.AddCommand(kinesisStreamAnalyserCmd)
}

var kinesisStreamAnalyserCmd = &cobra.Command{
	Use: "kinesis:stream:analyser",
	Short: "Analyse a Kinesis stream",
	Run: func(cmd *cobra.Command, args []string) {
		streamAnalyse()
	},
}

func streamAnalyse() {
	fmt.Println("AWS Kinesis Stream Analyser\n")
	fmt.Println("-----------------------------------------------------\n")

	stream, err := interact.AskForKinesisStream()
	if err != nil {
		fmt.Println(err)
		return
	}

	shardIds, keys := kinesis.GetShardIds(stream)

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

		recordCount := kinesis.GetShardRecordCount(shardIds[k], stream)
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