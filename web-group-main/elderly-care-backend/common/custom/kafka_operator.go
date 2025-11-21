package custom

import "github.com/segmentio/kafka-go"

type KafkaOperator struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}
