package memory

import (
	"sync"

	mailbox "github.com/xandalm/go-mailbox"
)

var (
	ErrEmptyBoxIdentifier    = mailbox.NewDetailedError(mailbox.ErrUnableToCreateBox, "identifier can't be empty")
	ErrRepeatedBoxIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToCreateBox, "repeated identifier")
)

type provider struct {
	mu    sync.RWMutex
	boxes map[string]*box
}

func NewProvider() *provider {
	return &provider{
		boxes: make(map[string]*box),
	}
}

func (p *provider) contains(id string) bool {
	_, ok := p.boxes[id]
	return ok
}

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	if id == "" {
		return nil, ErrEmptyBoxIdentifier
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.contains(id) {
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

func (p *provider) List() ([]string, mailbox.Error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	keys := []string{}
	for k := range p.boxes {
		keys = append(keys, k)
	}

	return keys, nil
}
