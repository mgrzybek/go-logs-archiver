package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"go-logs-archiver/internal/buffer"
	"go-logs-archiver/internal/core/domain"
	"go-logs-archiver/internal/core/ports"
	"go-logs-archiver/internal/lock"
	"go-logs-archiver/internal/producer"
)

const EXPECTED_MESSAGES = 1000

func configureLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()

	if err != nil {
		logger.Panic("Cannot create the logger")
	}

	return logger
}

func configureProducer(logger *zap.Logger) ports.MessagesProducer {
	result, err := producer.NewConsole()

	if err != nil {
		logger.Panic("Cannot create the buffer")
	}

	return result
}

func configureBuffer(logger *zap.Logger) ports.MessagesBuffer {
	result, err := buffer.NewMemoryBuffer(logger, 1)

	if err != nil {
		logger.Panic("Cannot create the buffer")
	}

	return result
}

func configureLock(logger *zap.Logger) ports.LockingSystem {
	result, err := lock.NewLocalMutex(logger, nil)

	if err != nil {
		logger.Panic("Cannot create the locking system")
	}

	return result
}

func generateMessages(logger *zap.Logger, engine *Engine, count int) {
	generatedMessages := 1
	logger.Sugar().Debugf("Generating %v the messages…", count)

	for generatedMessages <= count {
		logger.Sugar().Debugf("%v", generatedMessages)

		msg := domain.Message{
			Timestamp: time.Now().Unix(),
			Payload:   domain.RawMessage{},
		}
		logger.Sugar().Debug(msg)
		engine.ProcessMessage(msg)

		logger.Sugar().Debugf("%v messages generated", generatedMessages)
		generatedMessages++
	}

	logger.Debug("Generation ended")
}

func TestPurgeBuffer(t *testing.T) {
	logger := configureLogger()
	defer logger.Sync()

	logger.Debug("Creating the engine…")
	engine, err := NewEngine(
		logger,
		configureProducer(logger),
		configureBuffer(logger),
		configureLock(logger),
	)
	if err != nil {
		logger.Sugar().Panic(err)
	}
	logger.Debug("Engine created")

	generateMessages(logger, engine, EXPECTED_MESSAGES)
	engine.Terminate()

	assert.Equal(t, EXPECTED_MESSAGES, engine.metricMessagesFlushed)
}
