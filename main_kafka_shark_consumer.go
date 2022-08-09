package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/consumer"
	cluster "github.com/bsm/sarama-cluster"
)

func NewSharkConsumer() SharkTopic {
	ctx := "InitSharkConsumer"
	scoptTopic := "kafka_find_topic"
	topicSharkAccount := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_ACCOUNT")
	topicSharkAccountContact := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_ACCOUNT_CONTACT")
	topicSharkContact := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_CONTACT")
	topicSharkAddress := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_ADDRESS")
	topicSharkLeads := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_LEADS")
	topicSharkPhone := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_PHONE")
	topicSharkContactAddress := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_CONTACT_ADDRESS")
	topicSharkDocument := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_DOCUMENT")
	topicSharkContactDocument := golib.GetEnvOrFail(ctx, scoptTopic, "TOPIC_SHARK_CONTACT_DOCUMENT")

	return SharkTopic{
		Account:         topicSharkAccount,
		AccountContact:  topicSharkAccountContact,
		Contact:         topicSharkContact,
		Address:         topicSharkAddress,
		Leads:           topicSharkLeads,
		Phone:           topicSharkPhone,
		ContactAddress:  topicSharkContactAddress,
		Document:        topicSharkDocument,
		ContactDocument: topicSharkContactDocument,
	}
}

// for non-cdc synchronization between shark and sturgeon
func consumeKafkaShark(cfg localConfig.ServiceRepository, brokers []string) {
	config := cluster.NewConfig()
	config.ClientID = "sturgeon-shark-kafka"
	config.Group.Mode = cluster.ConsumerModePartitions
	topics := NewSharkConsumer()

	sharkTopics := []string{topics.Account, topics.AccountContact, topics.Contact, topics.Address,
		topics.Phone, topics.Document, topics.ContactDocument, topics.ContactAddress, topics.Leads}

	sharkConsumer, err := cluster.NewConsumer(brokers, "consumer-group-shark", sharkTopics, config)
	if err != nil {
		panic(err)
	}
	defer sharkConsumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case partition, ok := <-sharkConsumer.Partitions():
			if !ok {
				return
			}

			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					log.Printf("consuming %s offset %d", msg.Topic, msg.Offset)

					if err = consumeSharkTopics(cfg, msg.Topic, msg.Value, msg.Offset, topics); err != nil {
						helper.SendErrorLog(context.Background(), "Exec-Shark-Consumer", "general_parsing_shark", err, msg.Offset)
					}

					sharkConsumer.MarkOffset(msg, "")
					log.Printf("completed %s offset %d", msg.Topic, msg.Offset)
				}
			}(partition)

		case <-signals:
			os.Exit(1)
		}
	}
}

func consumeSharkTopics(cfg localConfig.ServiceRepository, topic string, payload []byte, offset int64, topics SharkTopic) (err error) {
	ctx := "consumeSharkTopics"
	tr := tracer.StartTrace(context.Background(), "consumeSharkTopics")
	tags := map[string]interface{}{}
	tags["topic"] = topic
	tags["offset"] = offset
	tags["message"] = string(payload)
	defer tr.Finish(tags)

	switch topic {
	// Topic table account
	case topics.Account:
		err = consumer.ExecSharkAccount(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table account contact
	case topics.AccountContact:
		err = consumer.ExecSharkAccountContact(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table contact
	case topics.Contact:
		err = consumer.ExecSharkContact(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table address
	case topics.Address:
		err = consumer.ExecSharkAddress(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table phone
	case topics.Phone:
		err = consumer.ExecSharkPhone(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table document
	case topics.Document:
		err = consumer.ExecSharkDocument(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table contact address
	case topics.ContactAddress:
		err = consumer.ExecSharkContactAddress(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table leads
	case topics.Leads:
		err = consumer.ExecSharkLeads(tr.NewChildContext(), cfg, ctx, payload)

		// Topic table contact document
	case topics.ContactDocument:
		err = consumer.ExecSharkContactDocument(tr.NewChildContext(), cfg, ctx, payload)
	}

	return err
}
