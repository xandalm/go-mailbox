package filesystem

import (
	"context"
	"os"
	"slices"
	"time"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrRepeatedContentIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "provided identifier is already in use")
	ErrPostingNilContent         = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "can't post nil content")
	ErrContentNotFound           = mailbox.NewDetailedError(mailbox.ErrUnableToReadContent, "not found")
)

type Bytes = mailbox.Bytes

func getContentIdentifiers(ch chan []string, f *os.File) []string {
	names, err := f.Readdirnames(0)
	if err != nil {
		ch <- nil
		return nil
	}
	ch <- names
	return names
}

type box struct {
	p  *provider
	bf *boxFile
}

// Clean implements mailbox.Box.
func (b *box) CleanWithContext(ctx context.Context) mailbox.Error {
	b.bf.mu.Lock()
	defer b.bf.mu.Unlock()

	f := b.bf.f

	var names []string

	ch1 := make(chan []string, 1)
	go getContentIdentifiers(ch1, f)
	select {
	case <-ctx.Done():
		return mailbox.ErrUnableToCleanBox
	case got := <-ch1:
		if got == nil {
			return mailbox.ErrUnableToCleanBox
		}
		names = got
	}
	errCount := 0
	for _, name := range names {
		if err := os.Remove(join(f.Name(), name)); err != nil {
			errCount++
		}
	}
	if errCount > 0 {
		return mailbox.ErrUnableToCleanBox
	}
	return nil
}

// Clean implements mailbox.Box.
func (b *box) Clean() mailbox.Error {
	return b.CleanWithContext(context.TODO())
}

// Delete implements mailbox.Box.
func (b *box) DeleteWithContext(_ context.Context, id string) mailbox.Error {
	b.bf.mu.Lock()
	defer b.bf.mu.Unlock()

	f := b.bf.f

	err := os.Remove(join(f.Name(), id))
	if err != nil && !os.IsNotExist(err) {
		return mailbox.ErrUnableToDeleteContent
	}
	return nil
}

// Delete implements mailbox.Box.
func (b *box) Delete(id string) mailbox.Error {
	return b.DeleteWithContext(context.TODO(), id)
}

// Get implements mailbox.Box.
func (b *box) GetWithContext(_ context.Context, id string) (mailbox.Data, mailbox.Error) {
	b.bf.mu.RLock()
	defer b.bf.mu.RUnlock()

	f := b.bf.f

	data, err := os.ReadFile(join(f.Name(), id))
	if err != nil {
		if os.IsNotExist(err) {
			return mailbox.Data{}, ErrContentNotFound
		}
		return mailbox.Data{}, mailbox.ErrUnableToReadContent
	}
	return mailbox.Data{
		Content: data,
	}, nil
}

// Get implements mailbox.Box.
func (b *box) Get(id string) (mailbox.Data, mailbox.Error) {
	return b.GetWithContext(context.TODO(), id)
}

// LazyGet implements mailbox.Box.
func (b *box) LazyGetWithContext(ctx context.Context, ids ...string) chan mailbox.AttemptData {

	ch := make(chan mailbox.AttemptData)

	go func() {
		get := func(ch chan mailbox.AttemptData, b *box, id string) chan struct{} {
			ch2 := make(chan struct{})

			go func() {
				data, err := b.Get(id)
				ch <- mailbox.AttemptData{
					Data:  data,
					Error: err,
				}
				close(ch2)
			}()

			return ch2
		}
		for i := 0; i < len(ids); i++ {
			select {
			case <-get(ch, b, ids[i]):
			case <-ctx.Done():
				i = len(ids)
			}
		}
		close(ch)
	}()

	return ch
}

// LazyGet implements mailbox.Box.
func (b *box) LazyGet(ids ...string) chan mailbox.AttemptData {
	return b.LazyGetWithContext(context.TODO(), ids...)
}

// ListFromPeriod implements mailbox.Box.
func (b *box) ListFromPeriodWithContext(ctx context.Context, begin, end time.Time, limit int) ([]string, mailbox.Error) {
	b.bf.mu.RLock()
	defer b.bf.mu.RUnlock()

	ret := make([]string, 0)
	idx := make([]int64, 0)

	f := b.bf.f

	path := f.Name()
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, mailbox.ErrUnableToReadContent
	}
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, mailbox.ErrUnableToReadContent
		}
		if info.ModTime().Before(begin) || info.ModTime().After(end) {
			continue
		}
		ct := info.ModTime().UnixNano()
		pos, _ := slices.BinarySearchFunc(idx, ct, func(a, b int64) int {
			return int(a - b)
		})
		idx = slices.Insert(idx, pos, ct)
		ret = slices.Insert(ret, pos, file.Name())
	}

	if limit <= 0 {
		return ret, nil
	}

	return ret[:limit], nil
}

// ListFromPeriod implements mailbox.Box.
func (b *box) ListFromPeriod(begin, end time.Time, limit int) ([]string, mailbox.Error) {
	return b.ListFromPeriodWithContext(context.TODO(), begin, end, limit)
}

// Post implements mailbox.Box.
func (b *box) PostWithContext(ctx context.Context, id string, c Bytes) (*time.Time, mailbox.Error) {
	b.bf.mu.Lock()
	defer b.bf.mu.Unlock()

	if c == nil {
		return nil, ErrPostingNilContent
	}

	f := b.bf.f

	name := join(f.Name(), id)
	if _, err := os.Stat(name); err == nil {
		return nil, ErrRepeatedContentIdentifier
	} else if !os.IsNotExist(err) {
		return nil, mailbox.ErrUnableToPostContent
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, mailbox.ErrUnableToPostContent
	}
	ct := time.Now()
	_, err = f.Write(c)
	if err != nil {
		f.Close()
		os.Remove(name)
		return nil, mailbox.ErrUnableToPostContent
	}
	if stat, err := f.Stat(); err == nil {
		ct = stat.ModTime()
	}
	f.Close()
	return &ct, nil
}

// Post implements mailbox.Box.
func (b *box) Post(id string, c Bytes) (*time.Time, mailbox.Error) {
	return b.PostWithContext(context.TODO(), id, c)
}
