package filesystem

import (
	"context"
	"io/fs"
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

func getDirectoryNamesFn[T []string](ch chan Data_Error[T], f *os.File) {
	names, err := f.Readdirnames(0)
	ch <- Data_Error[T]{names, err}
}

func getDirectoryNames(f *os.File) chan Data_Error[[]string] {
	ch := make(chan Data_Error[[]string], 1)
	go getDirectoryNamesFn(ch, f)
	return ch
}

func getDirectoryEntriesFn[T []fs.DirEntry](ch chan Data_Error[T], f *os.File) {
	entries, err := f.ReadDir(0)
	ch <- Data_Error[T]{entries, err}
}

func getDirectoryEntries(f *os.File) chan Data_Error[[]fs.DirEntry] {
	ch := make(chan Data_Error[[]fs.DirEntry], 1)
	go getDirectoryEntriesFn(ch, f)
	return ch
}

func getFileModTimeFn[T *time.Time](ch chan Data_Error[T], e fs.DirEntry) {
	info, err := e.Info()
	var data *time.Time
	if err == nil {
		t := info.ModTime()
		data = &t
	}
	ch <- Data_Error[T]{data, err}
}

func getFileModTime(e fs.DirEntry) chan Data_Error[*time.Time] {
	ch := make(chan Data_Error[*time.Time], 1)
	go getFileModTimeFn(ch, e)
	return ch
}

type Data_Error[T any] struct {
	data T
	err  error
}

func readFileContentFn[T []byte](ch chan Data_Error[T], name string) {
	data, err := os.ReadFile(name)
	ch <- Data_Error[T]{data, err}
}

func readFileContent(name string) chan Data_Error[[]byte] {
	ch := make(chan Data_Error[[]byte], 1)
	go readFileContentFn(ch, name)
	return ch
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

	select {
	case <-ctx.Done():
		return mailbox.ErrUnableToCleanBox
	case got := <-getDirectoryNames(f):
		if got.err != nil {
			return mailbox.ErrUnableToCleanBox
		}
		names = got.data
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

	got := <-readFileContent(join(f.Name(), id))
	if got.err != nil {
		if os.IsNotExist(got.err) {
			return mailbox.Data{}, ErrContentNotFound
		}
		return mailbox.Data{}, mailbox.ErrUnableToReadContent
	}
	return mailbox.Data{
		Content: got.data,
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

	var files []fs.DirEntry

	select {
	case <-ctx.Done():
		return ret, nil
	case got := <-getDirectoryEntries(f):
		if got.err != nil {
			return nil, mailbox.ErrUnableToReadContent
		}
		files = got.data
	}

	var err mailbox.Error
	for i := 0; i < len(files); i++ {
		file := files[i]
		select {
		case <-ctx.Done():
			i = len(files)
		case got := <-getFileModTime(file):
			if got.err != nil {
				err = mailbox.ErrUnableToReadContent
				i = len(files)
			} else if !(got.data.Before(begin) || got.data.After(end)) {
				ts := got.data.UnixNano()
				pos, _ := slices.BinarySearchFunc(idx, ts, func(a, b int64) int {
					return int(a - b)
				})
				idx = slices.Insert(idx, pos, ts)
				ret = slices.Insert(ret, pos, file.Name())
			}
		}
	}

	if limit <= 0 {
		return ret, err
	}

	return ret[:limit], err
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
