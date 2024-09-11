package mailbox

import (
	"slices"
	"sync"
)

type Error interface {
	sign() string
	Error() string
}

type mailboxError struct {
	msg string
}

func newError(msg string) *mailboxError {
	return &mailboxError{msg}
}

func (e mailboxError) sign() string {
	return "mailbox"
}

func (e mailboxError) Error() string {
	return e.sign() + ": " + e.msg
}

var (
	ErrBoxIDDuplicity Error = newError("duplicity of box identifier")
)

type Manager interface {
	// Create or restore a box.
	RequestBox(string) (Box, Error)
	// Remove a box and all its contents.
	EraseBox(string) Error
}

type Provider interface {
	// Create a new box.
	Create(string) (Box, Error)
	// Get existing box.
	Get(string) (Box, Error)
	// Delete existing box and all its contents.
	Delete(string) Error
	// List the identifier from all existing boxes.
	List() ([]string, Error)
}

type Box interface {
	// Post content.
	Post(any, any) Error
	// Read the corresponding content.
	Get(any) (any, Error)
	// Remove the corresponding content.
	Delete(any) Error
	// Remove all its existing contents.
	Clean() Error
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

func (m *manager) createBox(pos int, id string) (Box, Error) {
	b, err := m.p.Create(id)
	if err == nil {
		m.insert(pos, id)
	}
	return b, err
}

func (m *manager) getBox(id string) (Box, Error) {
	return m.p.Get(id)
}

func (m *manager) RequestBox(id string) (Box, Error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	has, pos := m.contains(id)
	if has {
		return m.getBox(id)
	}
	return m.createBox(pos, id)
}

func (m *manager) EraseBox(id string) Error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.p.Delete(id)
	return nil
}
