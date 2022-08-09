package service

import (
	"strconv"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

// 801	   1266845 ns/op	    4449 B/op	      76 allocs/op
// 996	   1213791 ns/op	    4375 B/op	      76 allocs/op
func BenchmarkKafkaConfig(b *testing.B) {
	pub, _ := NewKafkaPublisher("localhost:9092")
	for n := 0; n < b.N; n++ {
		pub.benchPub("test-benchmark", strconv.Itoa(n), []byte(`some message`))
	}
}

func (publisher *KafkaPublisherImpl) benchPub(topic string, messageKey string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(messageKey),
		Value:     sarama.ByteEncoder(message),
		Timestamp: time.Now(),
	}

	if _, _, err := publisher.producer.SendMessage(msg); err != nil {
		return err
	}
	defer publisher.producer.Close()

	return nil
}
