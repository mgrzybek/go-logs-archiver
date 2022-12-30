package core

// RawMessage is a simple array of bytes.
type RawMessage []byte

// RawMessages is an array of RawMessage
type RawMessages []RawMessage

// Message is a struct containing their extracted timestamp and their raw data
type Message struct {
	Timestamp int64 `json:"timestamp"`
	Payload   RawMessage
}

// MessagesConsumer gets the data from a broker
type MessagesConsumer interface {
	Run()
}

// MessagesProducer pushes the data to the persistent storage
type MessagesProducer interface {
	ProduceMessages(ts int64, messages RawMessages) (int, error)
}

// MessagesBuffer is used to store in a sorted way the messages while processing
type MessagesBuffer interface {
	PushMessage(message *Message) error
	PullAndDestroyMessages(ts int64) RawMessages

	GetTimestamps() []int64
}

// LockingSystem is used to lock the read/write processes of the other modules, using a local or a network-based tool
type LockingSystem interface {
	Lock(name string)
	Unlock()

	IsLocked() bool
}
