package memory

import (
	"container/list"
	"context"
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

func (b *box) CleanWithContext(_ context.Context) mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	clear(b.dataById)
	b.data = b.data.Init()
	return nil
}

func (b *box) Clean() mailbox.Error {
	return b.CleanWithContext(context.TODO())
}

func (b *box) DeleteWithContext(_ context.Context, k string) mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if e, ok := b.dataById[k]; ok {
		b.data.Remove(e)
		delete(b.dataById, k)
	}
	return nil
}

func (b *box) Delete(k string) mailbox.Error {
	return b.DeleteWithContext(context.TODO(), k)
}

func (b *box) GetWithContext(_ context.Context, k string) (mailbox.Data, mailbox.Error) {
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

func (b *box) Get(k string) (mailbox.Data, mailbox.Error) {
	return b.GetWithContext(context.TODO(), k)
}

func (b *box) LazyGetWithContext(ctx context.Context, ks ...string) chan mailbox.AttemptData {

	if ctx == nil {
		ctx = context.Background()
	}

	ch := make(chan mailbox.AttemptData)

	go func() {
		get := func(ch chan mailbox.AttemptData, b *box, k string) chan struct{} {
			ch2 := make(chan struct{})

			go func() {
				data, err := b.Get(k)
				ch <- mailbox.AttemptData{
					Data:  data,
					Error: err,
				}
				close(ch2)
			}()

			return ch2
		}
		for i := 0; i < len(ks); i++ {
			select {
			case <-get(ch, b, ks[i]):
			case <-ctx.Done():
				i = len(ks)
			}
		}
		close(ch)
	}()

	return ch
}

func (b *box) LazyGet(ks ...string) chan mailbox.AttemptData {
	return b.LazyGetWithContext(context.Background(), ks...)
}

func (b *box) ListFromPeriodWithContext(_ context.Context, begin, end time.Time, limit int) ([]string, mailbox.Error) {
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
	interrupt := func(s []string) bool { return false }
	if limit > 0 {
		interrupt = func(s []string) bool { return len(s) >= limit }
	}
	for {
		if interrupt(ids) || elem == nil {
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

func (b *box) ListFromPeriod(begin, end time.Time, limit int) ([]string, mailbox.Error) {
	return b.ListFromPeriodWithContext(context.TODO(), begin, end, limit)
}

func (b *box) PostWithContext(_ context.Context, id string, c Bytes) (*time.Time, mailbox.Error) {
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

func (b *box) Post(id string, c Bytes) (*time.Time, mailbox.Error) {
	return b.PostWithContext(context.TODO(), id, c)
}
