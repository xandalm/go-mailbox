package memory

import (
	"sync"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrPostingNilContent = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "can't post nil content")
)

type box struct {
	pk       uint64
	mu       sync.RWMutex
	contents map[any]any
}

func newBox() *box {
	return &box{
		contents: make(map[any]any),
	}
}

func (b *box) key() any {
	b.pk++
	return b.pk
}

func (b *box) Clean() mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	clear(b.contents)
	return nil
}

func (b *box) Delete(k any) mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.contents, k)
	return nil
}

func (b *box) Get(k any) (any, mailbox.Error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.contents[k], nil
}

func (b *box) Post(c any) (any, mailbox.Error) {
	if c == nil {
		return nil, ErrPostingNilContent
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	k := b.key()
	b.contents[k] = c
	return k, nil
}
