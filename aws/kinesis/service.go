package kinesis

import (
	"sort"
	"time"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/fatih/color"
)

var client *kinesis.Kinesis

func init() {
	client = _getClient()
}

func _getClient() *kinesis.Kinesis {

	region := "eu-west-1"
	if os.Getenv("AWS_DEFAULT_REGION") != "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}

	s := session.New(aws.NewConfig().WithRegion(region))
	s.Handlers.Send.PushFront(func(r *request.Request) {
		color.Set(color.FgYellow)
		//fmt.Printf("\nDEBUG: Request %s/%s\n\n", r.ClientInfo.ServiceName, r.Operation)
		color.Unset()
	})

	client := kinesis.New(s)
	return client
}

func GetStreams() (map[int]string, []int) {
	listStreamsOutput, err := client.ListStreams(nil)
	if err != nil {
		panic(err)
	}

	// TODO HasMoreStreams
	// fmt.Println(*listStreamsOutput.HasMoreStreams)

	streams := make(map[int]string)
	for key, value := range listStreamsOutput.StreamNames {
		streams[key] = *value
	}

	// The Go runtime actually randomizes the iteration order
	var keys []int
	for k := range streams {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return streams, keys
}

func GetShardIds(streamName string) ([]string, []int) {
	var streamIds []string
	describeStreamInput := &kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
	}
	client.DescribeStreamPages(describeStreamInput,
		func(stream *kinesis.DescribeStreamOutput, lastPage bool) bool {
			ids, _ := awsutil.ValuesAtPath(stream, "StreamDescription.Shards[].ShardId")
			for _, id := range ids {
				streamIds = append(streamIds, *id.(*string))
			}

			return !lastPage
		})

	// The Go runtime actually randomizes the iteration order
	var keys []int
	for k := range streamIds {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return streamIds, keys
}

func GetShardRecordCount(shardId string, streamName string) int {
	shardIteratorInput := &kinesis.GetShardIteratorInput{
		ShardId: aws.String(shardId),
		ShardIteratorType: aws.String("TRIM_HORIZON"),
		StreamName: aws.String(streamName),
	}
	getShardIteratorResult, _ := client.GetShardIterator(shardIteratorInput)
	shardIterator := *getShardIteratorResult.ShardIterator

	var recordsInShard int = 0
	for {
		recordsInput := &kinesis.GetRecordsInput{
			Limit: aws.Int64(10000),
			ShardIterator: aws.String(shardIterator),
		}
		result, _ := client.GetRecords(recordsInput)
		recordsInShard = recordsInShard + len(result.Records)
		shardIterator = *result.NextShardIterator

		time.Sleep(200 * 1000)
		if *result.MillisBehindLatest == 0 { break }
	}

	return recordsInShard
}