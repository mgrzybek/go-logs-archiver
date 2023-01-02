package ports

import (
	"go-logs-archiver/internal/core/domain"
)

// MessagesConsumer gets the data from a broker
type MessagesConsumer interface {
	Run()
}

// MessagesProducer pushes the data to the persistent storage
type MessagesProducer interface {
	ProduceMessages(ts int64, messages domain.RawMessages) (int, error)
}

// MessagesBuffer is used to store in a sorted way the messages while processing
type MessagesBuffer interface {
	PushMessage(message *domain.Message) error
	PullAndDestroyMessages(ts int64) domain.RawMessages

	GetTimestamps() []int64
}

// LockingSystem is used to lock the read/write processes of the other modules, using a local or a network-based tool
type LockingSystem interface {
	Lock(name string) error
	Unlock()

	IsLocked() bool
}
