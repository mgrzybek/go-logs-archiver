package consumer

import (
	"time"

	"go.uber.org/zap"

	"go-logs-archiver/internal/core"
	"go-logs-archiver/internal/core/domain"
)

// RandomGenerator is an object implementing interface MessagesConsumer
type RandomGenerator struct {
	engine           *core.Engine
	logger           *zap.Logger
	numberOfMessages int
}

// NewRandomGenerator is the constructor of the RandomGenerator
func NewRandomGenerator(logger *zap.Logger, engine *core.Engine, count int) (RandomGenerator, error) {
	return RandomGenerator{
		engine:           engine,
		logger:           logger,
		numberOfMessages: count,
	}, nil
}

func (r RandomGenerator) generateRandomMessage() domain.Message {
	return domain.Message{
		Timestamp: time.Now().Unix(),
		Payload:   domain.RawMessage{},
	}
}

// Run starts the consuming process
func (r RandomGenerator) Run() {
	generatedMessages := 0

	for generatedMessages <= r.numberOfMessages {
		r.engine.ProcessMessage(r.generateRandomMessage())
		generatedMessages++
	}
}
