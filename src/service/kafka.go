package service

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/Shopify/sarama"
)

//KafkaPublisherImpl struct
type KafkaPublisherImpl struct {
	producer sarama.SyncProducer
}

//NewKafkaPublisher constructor of PublisherImpl
func NewKafkaPublisher(addresses ...string) (*KafkaPublisherImpl, error) {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	configuration := sarama.NewConfig()
	configuration.Version = sarama.V0_11_0_0
	configuration.ClientID = "user-service-kafka"
	configuration.Producer.Retry.Max = 10
	configuration.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	configuration.Producer.RequiredAcks = sarama.WaitForLocal
	configuration.Producer.Timeout = 5 * time.Second
	configuration.Producer.Compression = sarama.CompressionSnappy
	configuration.Producer.Return.Successes = true

	// sync producer
	producer, err := sarama.NewSyncProducer(addresses, configuration)

	if err != nil {
		return nil, err
	}

	return &KafkaPublisherImpl{producer: producer}, nil
}

//Publish function
func (publisher *KafkaPublisherImpl) Publish(ctxReq context.Context, topic string, messageKey shared.MessageKey, message []byte) error {
	err := publisher.PublishKafka(ctxReq, topic, messageKey.String(), message)
	if err != nil {
		return err
	}
	return nil
}

//PublishKafka function
func (publisher *KafkaPublisherImpl) PublishKafka(ctxReq context.Context, topic string, messageKey string, message []byte) error {
	ctx := "KafkaService-PublishKafka"
	trace := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		trace.Finish(tags)
	}()

	// publish sync
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(messageKey),
		Value:     sarama.ByteEncoder(message),
		Timestamp: time.Now(),
	}

	tags["key"] = messageKey
	tags[helper.TextArgs] = msg

	_, _, err := publisher.producer.SendMessage(msg)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "publish_kafka", err, msg)
		return err
	}

	return nil
}

func (publisher *KafkaPublisherImpl) QueueJob(ctxReq context.Context, payload interface{}, messageKey, jobType string) error {
	ctx := "KafkaService-QueueJob"
	var auth string
	trace := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		trace.Finish(tags)
	}()
	topicWorker := golib.GetEnvOrFail(ctx, helper.TextFindServerKafkaConfig, "KAFKA_WORKER_TOPIC")

	// take token information from context
	authVal := ctxReq.Value(helper.TextAuthorization)
	if authVal != nil {
		auth = authVal.(string)
	}
	payloadM := serviceModel.QueuePayload{
		GeneralPayload: serviceModel.GeneralPayload{
			EventType: jobType,
			Payload:   payload,
		},
		Auth: auth,
	}
	byteMessage, err := json.Marshal(payloadM)
	if err != nil {
		return err
	}

	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topicWorker,
		Key:   sarama.StringEncoder(messageKey),
		Value: sarama.StringEncoder(byteMessage),
	}

	tags["key"] = messageKey
	tags[helper.TextArgs] = msg

	_, _, err = publisher.producer.SendMessage(msg)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "publish_kafka_worker", err, msg)
		return err
	}

	return nil
}

//BulkPublishKafka function
func (publisher *KafkaPublisherImpl) BulkPublishKafka(ctxReq context.Context, topic string, content []serviceModel.Messages) error {
	ctx := "KafkaService-PublishKafka"
	trace := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		trace.Finish(tags)
	}()

	messages := make([]*sarama.ProducerMessage, 0)
	for _, message := range content {
		// publish sync
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(message.Key),
			Value: sarama.StringEncoder(message.Content),
		}
		messages = append(messages, msg)
	}
	tags[helper.TextArgs] = messages

	if err := publisher.producer.SendMessages(messages); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "publish_kafka_batch", err, messages)
		return err
	}

	return nil
}
