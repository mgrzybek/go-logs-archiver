package lock

import (
	"sync"

	"go.uber.org/zap"
)

type LocalMutex struct {
	m      sync.Mutex
	logger *zap.Logger
}

func NewLockingSystem(logger *zap.Logger, ressourceName *string) (*LocalMutex, error) {
	return &LocalMutex{
		logger: logger,
	}, nil
}

func (l *LocalMutex) Lock() {
	l.m.Lock()
}

func (l *LocalMutex) Unlock() {
	l.m.Unlock()
}
