package domain

// RawMessage is a simple array of bytes.
type RawMessage []byte

// RawMessages is an array of RawMessage
type RawMessages []RawMessage

// Message is a struct containing their extracted timestamp and their raw data
type Message struct {
	Timestamp int64 `json:"timestamp"`
	Payload   RawMessage
}
