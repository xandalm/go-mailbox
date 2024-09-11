package mailbox

import (
	"errors"
	"slices"
)

var (
	ErrBoxIDDuplicity = errors.New("mailbox: id duplicity")
)

type Manager interface {
	CreateBox(string) (Box, error)
}

type Provider interface {
	Create(string) (Box, error)
	List() ([]string, error)
}

type Box interface {
}

type manager struct {
	p   Provider
	idx []string
}

func NewManager(f Provider) Manager {
	l, _ := f.List()
	m := &manager{f, l}
	return m
}

func (m *manager) contains(id string) (bool, int) {
	pos, ok := slices.BinarySearch(m.idx, id)
	return ok, pos
}

func (m *manager) insert(pos int, id string) {
	m.idx = slices.Insert(m.idx, pos, id)
}

func (m *manager) CreateBox(id string) (Box, error) {
	has, pos := m.contains(id)
	if has {
		return nil, ErrBoxIDDuplicity
	}
	b, _ := m.p.Create(id)
	m.insert(pos, id)
	return b, nil
}
