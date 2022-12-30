package producer

import (
	"fmt"

	"go-logs-archiver/internal/core"
)

// Console is an implementation of MessagesProducer
type Console struct {
	metricMessagesProducedTotal int
}

// NewConsole is the constructor of Console
func NewConsole() (Console, error) {
	return Console{
		metricMessagesProducedTotal: 0,
	}, nil
}

// ProduceMessages pushes the given messages into the persistent storage
func (c Console) ProduceMessages(ts int64, messages core.RawMessages) (int, error) {
	messageCounter := 0
	for _, m := range messages {
		fmt.Printf("%s\n", m)
		messageCounter++
	}
	c.metricMessagesProducedTotal += messageCounter

	return messageCounter, nil
}
