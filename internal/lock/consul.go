package lock

import (
	"go.uber.org/zap"

	"github.com/hashicorp/consul/api"
)

// ConsulLock is an implementation of LockingSystem
type ConsulLock struct {
	httpClient *api.Client
	logger     *zap.Logger
	lock       *api.Lock

	isLocked       bool
	brokenLockChan <-chan struct{}
}

// NewConsulLock is the constructor of ConsulLock
func NewConsulLock(logger *zap.Logger) (*ConsulLock, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		logger.Sugar().Error(err)
		return nil, err
	}

	return &ConsulLock{
		httpClient: client,
		isLocked:   false,
		logger:     logger,
	}, healthCheck(client, logger)
}

// Lock is used to get the lock
// https://github.com/hashicorp/consul/blob/main/command/lock/lock.go
func (c *ConsulLock) Lock(name string) error {
	var err error

	c.lock, err = c.httpClient.LockKey(name)
	if err != nil {
		c.logger.Sugar().Error(err)
		return err
	}

	c.brokenLockChan, err = c.lock.Lock(nil)
	if err != nil {
		c.logger.Sugar().Error(err)
		return err
	}

	c.isLocked = true
	go c.monitorLock()
	return nil
}

// Unlock is used to release the lock
func (c *ConsulLock) Unlock() {
	if err := c.lock.Unlock(); err != nil {
		c.logger.Sugar().Error(err)
	}
	c.isLocked = false
}

// IsLocked tells if the lock is still locked
func (c *ConsulLock) IsLocked() bool {
	return c.isLocked
}

func healthCheck(client *api.Client, logger *zap.Logger) error {
	leader, err := client.Status().Leader()
	if err != nil {
		logger.Sugar().Error(err)
		return err
	}

	logger.Sugar().Infof("Consul endpoint %s is healthy", leader)
	return nil
}

func (c *ConsulLock) monitorLock() {
	value := <-c.brokenLockChan
	c.logger.Sugar().Warn("The lock has been lost ", value)
	c.isLocked = false
}
