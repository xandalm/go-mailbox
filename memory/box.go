package memory

import (
	"sync"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrPostingNilContent         = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "can't post nil content")
	ErrRepeatedContentIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "provided identifier is already in use")
)

type Bytes = mailbox.Bytes

type box struct {
	mu       sync.RWMutex
	contents map[string]Bytes
}

func newBox() *box {
	return &box{
		contents: make(map[string]Bytes),
	}
}

func (b *box) Clean() mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	clear(b.contents)
	return nil
}

func (b *box) Delete(k string) mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.contents, k)
	return nil
}

func (b *box) Get(k string) (Bytes, mailbox.Error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.contents[k], nil
}

func (b *box) Post(id string, c Bytes) mailbox.Error {
	if c == nil {
		return ErrPostingNilContent
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.contents[id]; ok {
		return ErrRepeatedContentIdentifier
	}

	b.contents[id] = c
	return nil
}
