package core

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"go-logs-archiver/internal/core/domain"
	"go-logs-archiver/internal/core/ports"
)

// Engine is the main core object of the program.
type Engine struct {
	logger *zap.Logger

	producer   ports.MessagesProducer
	buffer     ports.MessagesBuffer
	flushMutex ports.LockingSystem

	wg sync.WaitGroup

	lastProcessedTimestamp int64
	timestampChannel       chan int64
	shouldTerminate        bool

	metricMessagesReceived int
	metricMessagesFlushed  int
}

// NewEngine is the constructor of the Engine.
func NewEngine(l *zap.Logger, p ports.MessagesProducer, b ports.MessagesBuffer, m ports.LockingSystem) (*Engine, error) {
	result := Engine{
		logger:                 l,
		producer:               p,
		buffer:                 b,
		flushMutex:             m,
		lastProcessedTimestamp: 0,
		timestampChannel:       make(chan int64),
		shouldTerminate:        false,
		metricMessagesReceived: 0,
		metricMessagesFlushed:  0,
	}

	go result.TriggerFlush()

	return &result, nil
}

// ProcessMessage receives a message from the consumer and pushes it into the buffer
func (e *Engine) ProcessMessage(m domain.Message) error {
	e.metricMessagesReceived++

	e.logger.Debug("Pushing a message…")
	if err := e.buffer.PushMessage(&m); err != nil {
		e.logger.Sugar().Error(err)
		return err
	}
	e.logger.Debug("Message pushed without error")

	if e.lastProcessedTimestamp < m.Timestamp {
		e.lastProcessedTimestamp = m.Timestamp
		e.timestampChannel <- e.lastProcessedTimestamp
		e.logger.Sugar().Debugf("New lastProcessedTimestamp set to %v", e.lastProcessedTimestamp)
	}
	return nil
}

// Terminate flushes the buffer before terminating the engine
func (e *Engine) Terminate() {
	e.logger.Info("Terminating the engine…")

	e.logger.Sugar().Debug("Closing the channel…")
	e.shouldTerminate = true
	close(e.timestampChannel)
	e.wg.Wait()

	e.logger.Info("Engine terminated, the buffer has been flushed.")
	e.logger.Sugar().Infof("Received messages: %v", e.metricMessagesReceived)
}

// flushBuffer is a private method used to regularly flush the buffer into the producer.
func (e *Engine) flushTimestampFromBuffer(ts int64) {
	count, err := e.producer.ProduceMessages(ts, e.buffer.PullAndDestroyMessages(ts))

	if err != nil {
		e.logger.Sugar().Error(err)
	} else {
		e.logger.Sugar().Debugf("%v messages pushed", count)
		e.metricMessagesFlushed += count
	}
}

// FlushBuffer is used to dump the whole buffer into the persistent storage and remove the data.
func (e *Engine) FlushBuffer() {
	for _, k := range e.buffer.GetTimestamps() {
		e.checkAndLock(k)
		e.flushTimestampFromBuffer(k)
	}
}

// TriggerFlush is a loop method watching the last timestamp inserted in order to flush the old data into the producer module.
// A mutex is used allow only one flushing process to run at a time in a clustered environmen for a given timestamp.
func (e *Engine) TriggerFlush() {
	e.logger.Info("Starting the TriggerFlush loop…")
	e.wg.Add(1)
	defer e.wg.Done()
	defer e.flushMutex.Unlock()

	for !e.shouldTerminate {
		newTimestamp := <-e.timestampChannel
		e.logger.Debug("Loop over the timestamps…")
		for _, k := range e.buffer.GetTimestamps() {
			if k < newTimestamp {
				e.checkAndLock(k)
				e.flushTimestampFromBuffer(k)
			}
		}
		e.logger.Debug("Loop over the timestamps ended")
	}

	e.logger.Debug("Terminating the TriggerFlush loop…")
	e.logger.Debug("Flush the remaining data from buffer…")
	e.FlushBuffer()

	e.logger.Debug("Flush ended")
	e.logger.Sugar().Infof("TriggerFlush loop terminated: %v messages processed", e.metricMessagesFlushed)
}

// checkAndLock locks the mutex if it has not been already done.
func (e *Engine) checkAndLock(ts int64) {
	if !e.flushMutex.IsLocked() {
		e.flushMutex.Lock(fmt.Sprintf("%c", ts))
		e.logger.Info("The lock is ours!")
	}
}
