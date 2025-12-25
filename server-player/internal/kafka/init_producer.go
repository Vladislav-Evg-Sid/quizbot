package kafka

import (
	"github.com/segmentio/kafka-go"
)

func NewProducer(brokers []string) Producer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}
