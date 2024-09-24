package mailbox

import "sync"

type memoryStorageBox struct{}

type MemoryStorage struct {
	mu    sync.RWMutex
	boxes map[string]memoryStorageBox
}

func (m *MemoryStorage) CreateBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.boxes[id]; ok {
		return ErrRepeatedBoxIdentifier
	}

	m.boxes[id] = memoryStorageBox{}
	return nil
}

func (m *MemoryStorage) ListBoxes() ([]string, Error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := []string{}
	for k := range m.boxes {
		ids = append(ids, k)
	}

	return ids, nil
}

func (m *MemoryStorage) DeleteBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.boxes, id)
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		boxes: make(map[string]memoryStorageBox),
	}
}
