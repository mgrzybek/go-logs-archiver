package lock

import (
	"sync"

	"go.uber.org/zap"
)

// LocalMutex is an implementation of LockingSystem
type LocalMutex struct {
	m      sync.Mutex
	logger *zap.Logger
}

// NewLockingSystem is the constructor of LocalMutex
func NewLockingSystem(logger *zap.Logger, ressourceName *string) (*LocalMutex, error) {
	return &LocalMutex{
		logger: logger,
	}, nil
}

// Lock is used to get the lock
func (l *LocalMutex) Lock(name string) {
	l.m.Lock()
}

// Unlock is used to release the lock
func (l *LocalMutex) Unlock() {
	l.m.Unlock()
}

// IsLocked tells that the lock is locked (always true)
func (l *LocalMutex) IsLocked() bool {
	return true
}
