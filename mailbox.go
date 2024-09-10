package mailbox

import (
	"errors"
	"slices"
)

var (
	ErrBoxIDDuplicity = errors.New("mailbox: id duplicity")
)

type Box struct {
	Id string
}

type Manager interface {
	CreateBox(string) (Box, error)
}

type Storage interface {
	Save(Box) error
	List() ([]string, error)
}

type manager struct {
	st  Storage
	idx []string
}

func NewManager(s Storage) Manager {
	l, _ := s.List()
	m := &manager{s, l}
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
		return Box{}, ErrBoxIDDuplicity
	}
	b := Box{id}
	m.st.Save(b)
	m.insert(pos, id)
	return b, nil
}
