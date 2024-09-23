package mailbox

import "sync"

type MemoryStorage struct {
	mu    sync.RWMutex
	boxes map[string]struct{}
}

func (m *MemoryStorage) CreateBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.boxes[id]; ok {
		return ErrRepeatedBoxIdentifier
	}

	m.boxes[id] = struct{}{}
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		boxes: make(map[string]struct{}),
	}
}
