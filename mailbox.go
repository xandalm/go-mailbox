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

func (e *mailboxError) sign() string {
	return "mailbox"
}

func (e *mailboxError) Error() string {
	return e.sign() + ": " + e.msg
}

type mailboxDetailedError struct {
	e    *mailboxError
	info string
}

// Returns a detailed error with the addicional info provided
func NewDetailedError(err Error, info string) Error {
	if err, ok := err.(*mailboxError); ok {
		return &mailboxDetailedError{err, info}
	}
	panic("invalid type of error")
}

func (e *mailboxDetailedError) sign() string {
	return e.e.Error()
}

func (e *mailboxDetailedError) Error() string {
	return e.sign() + ", " + e.info
}

var (
	ErrInvalidBoxIdentifier Error = newError("invalid box identifier")
	ErrUnknownBox           Error = newError("there's no such box")
	ErrInvalidPostContent   Error = newError("invalid content to post in box")
)

type Manager interface {
	// Create or restore a box.
	RequestBox(string) (Box, Error)
	// Remove box and all its contents.
	EraseBox(string) Error
	// Check if the box exists
	ContainsBox(string) bool
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
	// Post content and return its identifier.
	Post(any) (any, Error)
	// Read the content matching to the identifier.
	Get(any) (any, Error)
	// Remove the content matching to the identifier.
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
	has, pos := m.contains(id)
	if !has {
		return ErrUnknownBox
	}
	if err := m.p.Delete(id); err != nil {
		return err
	}
	m.idx = slices.Delete(m.idx, pos, pos+1)
	return nil
}

func (m *manager) ContainsBox(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	has, _ := m.contains(id)
	return has
}
