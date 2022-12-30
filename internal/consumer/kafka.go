package consumer

import (
	"encoding/json"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"go-logs-archiver/internal/core"
)

type Kafka struct {
	engine   *core.Engine
	logger   *zap.Logger
	worker   sarama.Consumer
	consumer sarama.PartitionConsumer
	topic    string

	shouldTerminate bool

	metricMessagesReceived int
}

func NewKafka(logger *zap.Logger, engine *core.Engine, bootstrapServers []string, groupID, offset, topic string) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	worker, err := sarama.NewConsumer(bootstrapServers, config)

	if err != nil {
		logger.Sugar().Error(err)
		return nil, err
	}

	return &Kafka{
		engine:          engine,
		logger:          logger,
		worker:          worker,
		topic:           topic,
		shouldTerminate: false,
	}, nil
}

func (k *Kafka) Run() {
	var err error

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	k.consumer, err = k.worker.ConsumePartition(k.topic, 0, sarama.OffsetOldest)
	if err != nil {
		k.logger.Sugar().Panic(err)
	}
	k.logger.Sugar().Infof("Subscribing to topic %v", k.topic)

	for !k.shouldTerminate {
		select {
		case err := <-k.consumer.Errors():
			k.logger.Sugar().Error(err)
		case msg := <-k.consumer.Messages():
			k.metricMessagesReceived++
			k.logger.Sugar().Debugf("Received message Count %d: | Topic(%s) | Message(%s)", k.metricMessagesReceived, string(msg.Topic), string(msg.Value))
			k.processMessage(msg.Value)
		case <-signals:
			k.logger.Info("SIGINT received, terminating…")
			k.shouldTerminate = true
		}
	}

	k.logger.Info("Closing the consumer…")
	err = k.consumer.Close()
	if err != nil {
		k.logger.Sugar().Error(err)
	}
	k.logger.Info("Closing the worker…")
	k.worker.Close()
	if err != nil {
		k.logger.Sugar().Error(err)
	}

	k.engine.Terminate()
}

func (k *Kafka) processMessage(message []byte) {
	buffer := core.Message{}
	err := json.Unmarshal(message, &buffer)
	if err != nil {
		k.logger.Sugar().Error(err)
		return
	}

	buffer.Payload = message
	k.logger.Sugar().Debugf("json object: %v", buffer)

	err = k.engine.ProcessMessage(buffer)
	if err != nil {
		k.logger.Sugar().Error(err)
	}
}
