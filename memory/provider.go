package memory

import (
	"context"
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

func (p *provider) Contains(ctx context.Context, id string) bool {
	_, ok := p.boxes[id]
	return ok
}

func (p *provider) Create(ctx context.Context, id string) (mailbox.Box, mailbox.Error) {
	if id == "" {
		return nil, ErrEmptyBoxIdentifier
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Contains(ctx, id) {
		return nil, ErrRepeatedBoxIdentifier
	}
	b := newBox()
	p.boxes[id] = b
	return b, nil
}

func (p *provider) Delete(ctx context.Context, id string) mailbox.Error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.boxes, id)
	return nil
}

func (p *provider) Get(ctx context.Context, id string) (mailbox.Box, mailbox.Error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if b, has := p.boxes[id]; has {
		return b, nil
	}
	return nil, mailbox.ErrBoxNotFound
}
