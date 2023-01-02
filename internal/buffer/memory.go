package buffer

import (
	"math/big"
	"sort"
	"sync"

	"go.uber.org/zap"

	"go-logs-archiver/internal/core/domain"
)

type storeMap map[int64]domain.RawMessages

// getKey provides the closest timestamp according to the given step
func getKey(ts, step int64) int64 {
	n, _ := big.NewInt(
		ts,
	).DivMod(
		big.NewInt(ts),
		big.NewInt(step),
		big.NewInt(0),
	)
	return n.Int64() * step
}

// MemoryBuffer is used to store in a sorted way the messages while processing
type MemoryBuffer struct {
	logger             *zap.Logger
	store              storeMap
	storageStepSeconds int64

	storageAccess sync.RWMutex

	metricStoreKeysSize int
}

// NewMemoryBuffer is the constructor of the in-memory map
func NewMemoryBuffer(logger *zap.Logger, secondsStep int64) (*MemoryBuffer, error) {
	return &MemoryBuffer{
		logger:             logger,
		store:              make(storeMap),
		storageStepSeconds: secondsStep,

		metricStoreKeysSize: 0,
	}, nil
}

// PushMessage inserts the given message into the map using the closest key
// TODO: sort the message using their own timestamp
func (b *MemoryBuffer) PushMessage(message *domain.Message) error {
	key := getKey(message.Timestamp, b.storageStepSeconds)
	b.logger.Sugar().Debugf("storing at key %v, value: %v", key, message.Payload)

	b.storageAccess.Lock()
	b.logger.Debug("Map locked")
	b.store[key] = append(b.store[key], message.Payload)
	b.metricStoreKeysSize = len(b.store)
	b.logger.Sugar().Debugf("Store unlocked: %v keys", b.metricStoreKeysSize)
	b.storageAccess.Unlock()

	return nil
}

// PullAndDestroyMessages returns the messages of the closest given key and destroys the data
func (b *MemoryBuffer) PullAndDestroyMessages(ts int64) domain.RawMessages {
	var result domain.RawMessages

	b.storageAccess.Lock()
	key := getKey(ts, b.storageStepSeconds)
	result = b.store[key]
	delete(b.store, key)
	b.storageAccess.Unlock()

	return result
}

// GetTimestamps provides a sorted list of the stored timestamps
func (b *MemoryBuffer) GetTimestamps() []int64 {
	b.storageAccess.RLock()
	b.logger.Debug("Map RLocked")

	result := make([]int64, 0, len(b.store))
	for k := range b.store {
		result = append(result, k)
	}
	b.storageAccess.RUnlock()
	b.logger.Debug("Map RUnlocked")

	sort.Slice(
		result,
		func(i, j int) bool { return result[i] < result[j] },
	)

	return result
}
