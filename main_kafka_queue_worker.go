package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	consumerExec "github.com/Bhinneka/user-service/src/consumer"
	cluster "github.com/bsm/sarama-cluster"
)

func dispatchWorker(appService *AppService, brokers []string) {
	ctx := "background_worker"
	scopeTopic := "kafka_find_topic"
	ctxReq := context.Background()

	topicWorker := golib.GetEnvOrFail(ctx, scopeTopic, "KAFKA_WORKER_TOPIC")
	config := cluster.NewConfig()
	config.ClientID = "user-service-worker"
	config.Group.Mode = cluster.ConsumerModePartitions

	// init consumer
	topics := []string{topicWorker}
	consumer, err := cluster.NewConsumer(brokers, "consumer-group", topics, config)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, scopeTopic, err, config)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case partition, ok := <-consumer.Partitions():
			if !ok {
				os.Exit(1)
			}

			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					tracer.WithTraceFunc(context.Background(), "SturgeonWorker", func(ctxReq context.Context, tags map[string]interface{}) {
						tags["topic"] = msg.Topic
						tags["partition"] = msg.Partition
						tags["offset"] = msg.Offset
						tags["message"] = string(msg.Value)

						switch msg.Topic {
						case topicWorker:
							consumerExec.Dispatch(ctxReq, msg.Value, appService.MerchantUseCase, appService.MemberUseCase, appService.ShippingAddressUseCase)
						}
						consumer.MarkOffset(msg, "")
					})

				}
			}(partition)

		case <-signals:
			os.Exit(1)
		}
	}
}
