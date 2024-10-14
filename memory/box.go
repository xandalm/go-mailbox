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

func (b *box) LazyGet(ks ...string) (chan mailbox.Data, mailbox.Error) {

	ch := make(chan mailbox.Data)

	go func() {

		var reg *registry

		for _, k := range ks {
			b.mu.Lock()

			if elem, ok := b.dataById[k]; !ok {
				ch <- mailbox.Data{}
			} else {
				reg = elem.Value.(*registry)
			}

			ch <- mailbox.Data{
				CreationTime: reg.ct,
				Content:      reg.c,
			}

			b.mu.Unlock()
		}

		close(ch)
	}()

	return ch, nil
}

func (b *box) ListFromPeriod(begin, end time.Time) ([]string, mailbox.Error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ids := make([]string, 0)

	elem := b.data.Front()
	for {
		if elem == nil {
			break
		}
		if elem.Value.(*registry).ct >= begin.UnixNano() {
			break
		}
		elem = elem.Next()
	}
	for {
		if elem == nil {
			break
		}
		reg := elem.Value.(*registry)
		if reg.ct > end.UnixNano() {
			break
		}
		ids = append(ids, reg.id)
		elem = elem.Next()
	}
	return ids, nil
}

func (b *box) Post(id string, c Bytes) (*time.Time, mailbox.Error) {
	if c == nil {
		return nil, ErrPostingNilContent
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.dataById[id]; ok {
		return nil, ErrRepeatedContentIdentifier
	}
	now := time.Now()
	reg := &registry{
		id,
		now.UnixNano(),
		c,
	}
	b.dataById[id] = b.data.PushBack(reg)
	return &now, nil
}
