package mailbox

import (
	"bytes"
	"sync"
)

type memoryStorageBox struct {
	content map[string][]byte
}

type MemoryStorage struct {
	mu    sync.RWMutex
	boxes map[string]*memoryStorageBox
}

func (m *MemoryStorage) CreateBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.boxes[id]; ok {
		return ErrRepeatedBoxIdentifier
	}

	m.boxes[id] = &memoryStorageBox{}
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

func (m *MemoryStorage) CleanBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if box, ok := m.boxes[id]; ok {
		clear(box.content)
	}
	return nil
}

func (m *MemoryStorage) CreateContent(bid, cid string, c []byte) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	box, ok := m.boxes[bid]
	if !ok {
		return ErrBoxNotFoundToPostContent
	}
	if box.content == nil {
		box.content = make(map[string][]byte)
	}
	if _, ok := box.content[cid]; ok {
		return ErrRepeatedContentIdentifier
	}
	box.content[cid] = bytes.Clone(c)
	return nil
}

func (m *MemoryStorage) ReadContent(bid, cid string) ([]byte, Error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	box, ok := m.boxes[bid]
	if !ok {
		return nil, ErrBoxNotFoundToReadContent
	}

	return box.content[cid], nil
}

func (m *MemoryStorage) DeleteContent(bid, cid string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	box, ok := m.boxes[bid]
	if !ok {
		return nil
	}
	delete(box.content, cid)
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		boxes: make(map[string]*memoryStorageBox),
	}
}
