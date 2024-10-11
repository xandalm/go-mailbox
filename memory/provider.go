package memory

import (
	"sync"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrEmptyBoxIdentifier    = mailbox.NewDetailedError(mailbox.ErrUnableToCreateBox, "identifier can't be empty")
	ErrRepeatedBoxIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToCreateBox, "repeated identifier")
)

type provider struct {
	mu    sync.RWMutex
	boxes map[string]*box
}

func NewProvider() mailbox.Provider {
	return &provider{
		boxes: make(map[string]*box),
	}
}

func (p *provider) Contains(id string) (bool, mailbox.Error) {
	_, ok := p.boxes[id]
	return ok, nil
}

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	if id == "" {
		return nil, ErrEmptyBoxIdentifier
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if has, _ := p.Contains(id); has {
		return nil, ErrRepeatedBoxIdentifier
	}
	b := newBox()
	p.boxes[id] = b
	return b, nil
}

func (p *provider) Delete(id string) mailbox.Error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.boxes, id)
	return nil
}

func (p *provider) Get(id string) (mailbox.Box, mailbox.Error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	b := p.boxes[id]
	return b, nil
}
