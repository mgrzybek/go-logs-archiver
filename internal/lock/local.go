package lock

import (
	"sync"

	"go.uber.org/zap"
)

// LocalMutex is an implementation of LockingSystem
type LocalMutex struct {
	locked bool
	m      sync.Mutex
	logger *zap.Logger
}

// NewLocalMutex is the constructor of LocalMutex
func NewLocalMutex(logger *zap.Logger, ressourceName *string) (*LocalMutex, error) {
	return &LocalMutex{
		locked: false,
		logger: logger,
	}, nil
}

// Lock is used to get the lock
func (l *LocalMutex) Lock(name string) error {
	l.m.Lock()
	l.locked = true
	return nil
}

// Unlock is used to release the lock
func (l *LocalMutex) Unlock() {
	l.m.Unlock()
	l.locked = false
}

// IsLocked tells that the lock is locked
func (l *LocalMutex) IsLocked() bool {
	return l.locked
}
