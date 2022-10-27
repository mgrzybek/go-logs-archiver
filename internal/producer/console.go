package producer

import (
	"fmt"

	"go-logs-archiver/internal/core"
)

type Console struct{
	metricMessagesProducedTotal int
}

func NewConsole() (Console, error) {
	return Console{
		metricMessagesProducedTotal: 0,
	}, nil
}

func (c Console) ProduceMessages(ts int64, messages core.RawMessages) (int, error) {
	messageCounter := 0
	for _, m := range messages {
		fmt.Printf("%s\n", m)
		messageCounter++
	}
	c.metricMessagesProducedTotal += messageCounter

	return messageCounter, nil
}
