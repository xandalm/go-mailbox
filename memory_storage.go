package mailbox

import (
	"bytes"
	"sync"
)

type memoryStorageBox struct {
	content map[string][]byte
}

type memoryStorage struct {
	mu    sync.RWMutex
	boxes map[string]*memoryStorageBox
}

func (m *memoryStorage) CreateBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.boxes[id]; ok {
		return ErrRepeatedBoxIdentifier
	}

	m.boxes[id] = &memoryStorageBox{
		content: make(map[string][]byte),
	}
	return nil
}

func (m *memoryStorage) ListBoxes() ([]string, Error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := []string{}
	for k := range m.boxes {
		ids = append(ids, k)
	}

	return ids, nil
}

func (m *memoryStorage) DeleteBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.boxes, id)
	return nil
}

func (m *memoryStorage) CleanBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if box, ok := m.boxes[id]; ok {
		clear(box.content)
	}
	return nil
}

func (m *memoryStorage) CreateContent(bid, cid string, c []byte) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	box, ok := m.boxes[bid]
	if !ok {
		return ErrBoxNotFoundToPostContent
	}
	if _, ok := box.content[cid]; ok {
		return ErrRepeatedContentIdentifier
	}
	box.content[cid] = bytes.Clone(c)
	return nil
}

func (m *memoryStorage) ReadContent(bid, cid string) ([]byte, Error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	box, ok := m.boxes[bid]
	if !ok {
		return nil, ErrBoxNotFoundToReadContent
	}

	return box.content[cid], nil
}

func (m *memoryStorage) DeleteContent(bid, cid string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()

	box, ok := m.boxes[bid]
	if !ok {
		return nil
	}
	delete(box.content, cid)
	return nil
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		boxes: make(map[string]*memoryStorageBox),
	}
}