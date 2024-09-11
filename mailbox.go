package mailbox

import (
	"errors"
	"slices"
	"sync"
)

var (
	ErrBoxIDDuplicity = errors.New("mailbox: id duplicity")
)

type Manager interface {
	RequestBox(string) (Box, error)
	EraseBox(string) error
}

type Provider interface {
	Create(string) (Box, error)
	Get(string) (Box, error)
	Delete(string) error
	List() ([]string, error)
}

type Box interface {
}

type manager struct {
	mu  sync.Mutex
	p   Provider
	idx []string
}

func NewManager(p Provider) Manager {
	l, _ := p.List()
	m := &manager{p: p, idx: l}
	return m
}

func (m *manager) contains(id string) (bool, int) {
	pos, ok := slices.BinarySearch(m.idx, id)
	return ok, pos
}

func (m *manager) insert(pos int, id string) {
	m.idx = slices.Insert(m.idx, pos, id)
}

func (m *manager) createBox(pos int, id string) (Box, error) {
	b, err := m.p.Create(id)
	if err == nil {
		m.insert(pos, id)
	}
	return b, err
}

func (m *manager) getBox(id string) (Box, error) {
	return m.p.Get(id)
}

func (m *manager) RequestBox(id string) (Box, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	has, pos := m.contains(id)
	if has {
		return m.getBox(id)
	}
	return m.createBox(pos, id)
}

func (m *manager) EraseBox(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.p.Delete(id)
	return nil
}
