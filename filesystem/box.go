package filesystem

import (
	"os"
	"slices"
	"sync"
	"time"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrRepeatedContentIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "provided identifier is already in use")
	ErrPostingNilContent         = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "can't post nil content")
	ErrContentNotFound           = mailbox.NewDetailedError(mailbox.ErrUnableToReadContent, "not found")
)

type Bytes = mailbox.Bytes

type box struct {
	mu sync.RWMutex
	f  *os.File
	p  *provider
	id string
}

// Clean implements mailbox.Box.
func (b *box) Clean() mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	names, err := b.f.Readdirnames(0)
	if err != nil {
		return mailbox.ErrUnableToCleanBox
	}
	for _, name := range names {
		os.Remove(join(b.f.Name(), name))
	}
	return nil
}

// Delete implements mailbox.Box.
func (b *box) Delete(id string) mailbox.Error {
	b.mu.Lock()
	defer b.mu.Unlock()

	name := join(b.f.Name(), id)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return mailbox.ErrUnableToReadContent
	}
	if err := os.Remove(name); err != nil {
		return mailbox.ErrUnableToDeleteContent
	}
	return nil
}

// Get implements mailbox.Box.
func (b *box) Get(id string) (mailbox.Data, mailbox.Error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	name := join(b.f.Name(), id)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return mailbox.Data{}, ErrContentNotFound
	} else if err != nil {
		return mailbox.Data{}, mailbox.ErrUnableToReadContent
	}
	data, err := os.ReadFile(name)
	if err != nil {
		return mailbox.Data{}, mailbox.ErrUnableToReadContent
	}
	return mailbox.Data{
		Content: data,
	}, nil
}

// Get implements mailbox.Box.
func (b *box) LazyGet(ids ...string) chan mailbox.AttemptData {

	ch := make(chan mailbox.AttemptData)

	go func() {
		for _, id := range ids {
			data, err := b.Get(id)
			ch <- mailbox.AttemptData{
				Data:  data,
				Error: err,
			}
		}
		close(ch)
	}()

	return ch
}

// ListFromPeriod implements mailbox.Box.
func (b *box) ListFromPeriod(begin, end time.Time) ([]string, mailbox.Error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ret := make([]string, 0)
	idx := make([]int64, 0)

	path := b.f.Name()
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

	return ret, nil
}

// Post implements mailbox.Box.
func (b *box) Post(id string, c Bytes) (*time.Time, mailbox.Error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if c == nil {
		return nil, ErrPostingNilContent
	}
	name := join(b.f.Name(), id)
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
