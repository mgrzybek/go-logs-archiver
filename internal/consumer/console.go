package consumer

import (
	"bufio"
	"encoding/json"
	"os"

	"go.uber.org/zap"

	"go-logs-archiver/internal/core"
)

// Console is an object implementing interface MessagesConsumer
type Console struct {
	engine  *core.Engine
	logger  *zap.Logger
	scanner *bufio.Scanner
}

// NewConsole is the constructor of the Console
func NewConsole(logger *zap.Logger, engine *core.Engine) (Console, error) {
	return Console{
		engine:  engine,
		logger:  logger,
		scanner: bufio.NewScanner(os.Stdin),
	}, nil
}

// Run starts the consuming process
func (c Console) Run() {
	for c.scanner.Scan() {
		c.logger.Sugar().Debugf("received: %v", c.scanner.Text())

		buffer := core.Message{}
		err := json.Unmarshal(c.scanner.Bytes(), &buffer)

		if err != nil {
			c.logger.Sugar().Error(err)
		}
		buffer.Payload = c.scanner.Bytes()
		c.logger.Sugar().Debugf("json object: %v", buffer)

		c.engine.ProcessMessage(buffer)
	}

	if err := c.scanner.Err(); err != nil {
		c.logger.Error(err.Error())
	}
}
