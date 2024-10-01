package memory

import (
	"container/list"
	"sync"
	"time"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrPostingNilContent         = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "can't post nil content")
	ErrRepeatedContentIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "provided identifier is already in use")
)

type Bytes = mailbox.Bytes

type registry struct {
	id string
	ct int64
	c  Bytes
}

type box struct {
	mu       sync.RWMutex
	data     *list.List
	dataById map[string]*list.Element
}

func newBox() *box {
	return &box{
		data:     list.New(),
		dataById: map[string]*list.Element{},
	}
}

func (b *box) Clean() mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	clear(b.dataById)
	b.data = b.data.Init()
	return nil
}

func (b *box) Delete(k string) mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if e, ok := b.dataById[k]; ok {
		b.data.Remove(e)
		delete(b.dataById, k)
	}
	return nil
}

func (b *box) Get(k string) (mailbox.Data, mailbox.Error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var reg *registry
	if elem, ok := b.dataById[k]; !ok {
		return mailbox.Data{}, nil
	} else {
		reg = elem.Value.(*registry)
	}

	data := mailbox.Data{
		CreationTime: reg.ct,
		Content:      reg.c,
	}

	return data, nil
}

func (b *box) GetFromPeriod(int64, int64) ([]mailbox.Data, mailbox.Error) {
	panic("unimplemented")
}

func (b *box) Post(id string, c Bytes) (int64, mailbox.Error) {
	if c == nil {
		return 0, ErrPostingNilContent
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.dataById[id]; ok {
		return 0, ErrRepeatedContentIdentifier
	}
	reg := &registry{
		id,
		time.Now().UnixNano(),
		c,
	}
	b.dataById[id] = b.data.PushBack(reg)
	return reg.ct, nil
}
