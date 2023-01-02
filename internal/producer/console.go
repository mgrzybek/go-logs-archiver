package producer

import (
	"fmt"

	"go-logs-archiver/internal/core/domain"
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
func (c Console) ProduceMessages(ts int64, messages domain.RawMessages) (int, error) {
	messageCounter := 0
	for _, m := range messages {
		_, err := fmt.Printf("%s\n", m)
		if err != nil {
			return c.metricMessagesProducedTotal, err
		}
		messageCounter++
	}
	c.metricMessagesProducedTotal += messageCounter

	return messageCounter, nil
}
