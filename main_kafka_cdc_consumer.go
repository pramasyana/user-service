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

type SharkTopic struct {
	Account          string
	AccountContact   string
	AccountTemporary string
	Contact          string
	Address          string
	Phone            string
	Document         string
	ContactNPWP      string
	ContactAddress   string
	ContactTemp      string
	Leads            string
	Merchant         string
	MerchantDocument string
	MerchantBank     string
	ShippingAddress  string
	ContactDocument  string
}

func NewCDCConsumer() SharkTopic {
	ctx := "InitCDC"
	scoptTopic := "kafka_find_topic"
	topicSharkAccount := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_ACCOUNT")
	topicSharkAccountContact := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_ACCOUNT_CONTACT")
	topicSharkAccountTemporary := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_ACCOUNT_TEMPORARY")
	topicSharkContact := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_CONTACT")
	topicSharkAddress := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_ADDRESS")
	topicSharkPhone := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_PHONE")
	topicSharkDocument := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_DOCUMENT")
	topicSharkContactNpwp := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_CONTACT_NPWP")
	topicSharkContactAddress := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_CONTACT_ADDRESS")
	topicSharkContactTemp := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_CONTACT_TEMP")
	topicSharkLeads := golib.GetEnvOrFail(ctx, scoptTopic, "KAFKA_SHARK_LEADS")

	return SharkTopic{
		Account:          topicSharkAccount,
		AccountContact:   topicSharkAccountContact,
		AccountTemporary: topicSharkAccountTemporary,
		Contact:          topicSharkContact,
		Address:          topicSharkAddress,
		Phone:            topicSharkPhone,
		Document:         topicSharkDocument,
		ContactNPWP:      topicSharkContactNpwp,
		ContactAddress:   topicSharkContactAddress,
		ContactTemp:      topicSharkContactTemp,
		Leads:            topicSharkLeads,
	}
}

func consumeKafkaCDC(cfg localConfig.ServiceRepository) {
	ctx := "consume_kafka"

	kafkaBroker1 := golib.GetEnvOrFail(ctx, "kafka_find_host_1", "KAFKA_BHINNEKA_CDC_BROKER")

	// cluster kafka construct with partitions mode
	configCDC := cluster.NewConfig()
	configCDC.ClientID = "user-service-cdc-kafka"
	configCDC.Group.Mode = cluster.ConsumerModePartitions
	topics := NewCDCConsumer()

	// init consumerCDC
	brokersCDC := []string{kafkaBroker1}
	topicsCDC := []string{topics.Account, topics.AccountContact, topics.AccountTemporary, topics.Contact, topics.Address,
		topics.Phone, topics.Document, topics.ContactNPWP, topics.ContactAddress, topics.ContactTemp, topics.Leads}

	consumerCDC, err := cluster.NewConsumer(brokersCDC, "consumer-group-cdc", topicsCDC, configCDC)
	if err != nil {
		os.Exit(1)
	}
	defer consumerCDC.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case partition, ok := <-consumerCDC.Partitions():
			if !ok {
				return
			}

			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					log.Printf("consuming %s offset %d", msg.Topic, msg.Offset)

					if err = consumeCDC(cfg, msg.Topic, msg.Value, msg.Offset, topics); err != nil {
						helper.SendErrorLog(context.Background(), "Exec-CDC", "general_parsing_cdc", err, msg.Offset)
					}

					consumerCDC.MarkOffset(msg, "")
					log.Printf("completed %s offset %d", msg.Topic, msg.Offset)
				}
			}(partition)

		case <-signals:
			os.Exit(1)
		}
	}
}

func consumeCDC(cfg localConfig.ServiceRepository, topic string, payload []byte, offset int64, topics SharkTopic) (err error) {
	ctx := "consumeCDC"
	tr := tracer.StartTrace(context.Background(), "KafkaCDCConsumer")
	tags := map[string]interface{}{}
	tags["topic"] = topic
	tags["offset"] = offset
	tags["message"] = string(payload)
	defer tr.Finish(tags)

	switch topic {
	// Topic table account
	case topics.Account:
		err = consumer.ProcessSharkAccount(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table account temporary
	case topics.AccountTemporary:
		err = consumer.ProcessSharkAccountTemp(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table account contact
	case topics.AccountContact:
		err = consumer.ProcessSharkAccountContact(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table contact
	case topics.Contact:
		err = consumer.ProcessSharkContact(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table address
	case topics.Address:
		err = consumer.ProcessSharkAddress(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table phone
	case topics.Phone:
		err = consumer.ProcessSharkPhone(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table document
	case topics.Document:
		err = consumer.ProcessSharkDocument(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table contact npwp
	case topics.ContactNPWP:
		err = consumer.ProcessSharkContactNPWP(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table contact address
	case topics.ContactAddress:
		err = consumer.ProcessSharkContactAddress(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table contact temp
	case topics.ContactTemp:
		err = consumer.ProcessSharkContactTemp(tr.NewChildContext(), cfg, ctx, payload)

	// Topic table leads
	case topics.Leads:
		err = consumer.ProcessSharkLeads(tr.NewChildContext(), cfg, ctx, payload)

	}

	return err
}
