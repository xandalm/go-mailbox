package mailbox

import "sync"

type MemoryStorage struct {
	mu    sync.RWMutex
	boxes map[string]struct{}
}

func (m *MemoryStorage) CreateBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.boxes[id] = struct{}{}
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		boxes: make(map[string]struct{}),
	}
}
