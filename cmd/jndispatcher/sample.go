package main

import (
    "context"
    "fmt"
    "log"

    "github.com/twmb/franz-go/pkg/kgo"
)

func main() {
    client, err := kgo.NewClient(
        kgo.SeedBrokers("localhost:9092"),
        kgo.ConsumerGroup("my-consumer-group"),
        kgo.ConsumeTopics("topic1", "topic4"),
        kgo.ClientID("sample-consumer"),
    )
    if err != nil {
        log.Fatalf("unable to create client: %v", err)
    }
    defer client.Close()

    for {
        fetches := client.PollFetches(context.Background())
        if errs := fetches.Errors(); len(errs) > 0 {
            for _, fetchErr := range errs {
                log.Printf("fetch error: %v", fetchErr)
            }
            continue
        }

        iter := fetches.RecordIter()
        for !iter.Done() {
            record := iter.Next()
            handleKafkaMessage(record)
        }
    }
}

// handleKafkaMessage simulates processing like a web request.
func handleKafkaMessage(record *kgo.Record) {
    fmt.Printf("Topic: %s Partition: %d Offset: %d Value: %s\n",
        record.Topic, record.Partition, record.Offset, string(record.Value))
    // Here goes the "web service" style handling logic.
}
