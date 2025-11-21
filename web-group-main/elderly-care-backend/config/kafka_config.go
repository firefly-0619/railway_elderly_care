package config

import (
	"elderly-care-backend/common/constants"
	"elderly-care-backend/common/custom"
	. "elderly-care-backend/global"
	"github.com/segmentio/kafka-go"
)

func initKafka() {
	conn, err := kafka.Dial("tcp", Config.Kafka.Addresses[0])
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	kafkaTopics := []kafka.TopicConfig{
		{
			Topic:             constants.MESSAGE_TOPIC,
			NumPartitions:     2,
			ReplicationFactor: 1,
		},
	}
	err = conn.CreateTopics(kafkaTopics...)
	if err != nil {
		panic(err)
	}
	Logger.Info("kafka client init success")
	kafkaOperatorMap := make(map[string]*custom.KafkaOperator)
	// 创建kafka生产者和消费者
	kafkaOperatorMap[constants.MESSAGE_TOPIC] = &custom.KafkaOperator{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(Config.Kafka.Addresses...),
			Topic:        constants.MESSAGE_TOPIC,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireNone,
			Async:        Config.Kafka.Async,
		},
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: Config.Kafka.Addresses,
			Topic:   constants.MESSAGE_TOPIC,
			GroupID: constants.MESSAGE_GROUP,
		}),
	}

	KafkaOperators = kafkaOperatorMap

}
