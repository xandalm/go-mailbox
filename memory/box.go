package memory

import (
	"sync"

	"github.com/xandalm/go-mailbox"
)

type box struct {
	pk       uint64
	mu       sync.RWMutex
	contents map[any]any
}

func (b *box) key() any {
	k := b.pk
	b.pk++
	return k
}

func (b *box) Clean() mailbox.Error {
	panic("unimplemented")
}

func (b *box) Delete(any) mailbox.Error {
	panic("unimplemented")
}

func (b *box) Get(k any) (any, mailbox.Error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.contents[k], nil
}

func (b *box) Post(c any) (any, mailbox.Error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	k := b.key()
	b.contents[k] = c
	return k, nil
}
