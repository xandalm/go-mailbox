package memory

import (
	"sync"

	mailbox "github.com/xandalm/go-mailbox"
)

type provider struct {
	mu    sync.RWMutex
	boxes map[string]*box
}

func (p *provider) contains(id string) bool {
	_, ok := p.boxes[id]
	return ok
}

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.contains(id) {
		return nil, mailbox.ErrBoxIDDuplicity
	}
	b := &box{}
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
