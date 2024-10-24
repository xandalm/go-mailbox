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

type ioResult[T any] struct {
	data T
	err  error
}

func getDirectoryNamesFn[T []string](ch chan ioResult[T], f *os.File) {
	f.Seek(0, 0)
	names, err := f.Readdirnames(0)
	ch <- ioResult[T]{names, err}
}

func getDirectoryNames(f *os.File) chan ioResult[[]string] {
	ch := make(chan ioResult[[]string], 1)
	go getDirectoryNamesFn(ch, f)
	return ch
}

func getDirectoryEntriesFn[T []fs.DirEntry](ch chan ioResult[T], f *os.File) {
	f.Seek(0, 0)
	entries, err := f.ReadDir(0)
	ch <- ioResult[T]{entries, err}
}

func getDirectoryEntries(f *os.File) chan ioResult[[]fs.DirEntry] {
	ch := make(chan ioResult[[]fs.DirEntry], 1)
	go getDirectoryEntriesFn(ch, f)
	return ch
}

func getFileInfoFromDirEntryFn[T *fs.FileInfo](ch chan ioResult[T], e fs.DirEntry) {
	info, err := e.Info()
	ch <- ioResult[T]{&info, err}
}

func getFileInfoFromDirEntry(e fs.DirEntry) chan ioResult[*fs.FileInfo] {
	ch := make(chan ioResult[*fs.FileInfo], 1)
	go getFileInfoFromDirEntryFn(ch, e)
	return ch
}

func getFileInfoFn[T *fs.FileInfo](ch chan ioResult[T], name string) {
	info, err := os.Stat(name)
	ch <- ioResult[T]{&info, err}
}

func getFileInfo(name string) chan ioResult[*fs.FileInfo] {
	ch := make(chan ioResult[*fs.FileInfo], 1)
	go getFileInfoFn(ch, name)
	return ch
}

func readFileContentFn[T []byte](ch chan ioResult[T], name string) {
	data, err := os.ReadFile(name)
	ch <- ioResult[T]{data, err}
}

func readFileContent(name string) chan ioResult[[]byte] {
	ch := make(chan ioResult[[]byte], 1)
	go readFileContentFn(ch, name)
	return ch
}

func openFileFn[T *os.File](ch chan ioResult[T], name string) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	ch <- ioResult[T]{f, err}
}

func openFile(name string) chan ioResult[*os.File] {
	ch := make(chan ioResult[*os.File], 1)
	go openFileFn(ch, name)
	return ch
}

func writeContentFn(ch chan bool, f *os.File, c []byte) {
	_, err := f.Write(c)
	ch <- (err == nil)
}

func writeContent(f *os.File, c []byte) chan bool {
	ch := make(chan bool, 1)
	go writeContentFn(ch, f, c)
	return ch
}

type box struct {
	p  *provider
	bf *boxFile
}

// CleanWithContext implements mailbox.Box.
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

// DeleteWithContext implements mailbox.Box.
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

// GetWithContext implements mailbox.Box.
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

// LazyGetWithContext implements mailbox.Box.
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

// ListFromPeriodWithContext implements mailbox.Box.
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

	nfiles := len(files)

	if limit <= 0 {
		limit = nfiles
	}

	var err mailbox.Error
	for i := 0; i < nfiles; i++ {
		file := files[i]
		select {
		case <-ctx.Done():
			i = nfiles
		case got := <-getFileInfoFromDirEntry(file):
			if got.err != nil {
				err = mailbox.ErrUnableToReadContent
				i = nfiles
			} else {
				ct := (*got.data).ModTime()
				if !(ct.Before(begin) || ct.After(end)) {
					ts := ct.UnixNano()
					pos, _ := slices.BinarySearchFunc(idx, ts, func(a, b int64) int {
						return int(a - b)
					})
					idx = slices.Insert(idx, pos, ts)
					ret = slices.Insert(ret, pos, file.Name())
				}
			}
		}
	}

	if lim := len(ret); lim < limit {
		limit = lim
	}

	return ret[:limit], err
}

// ListFromPeriod implements mailbox.Box.
func (b *box) ListFromPeriod(begin, end time.Time, limit int) ([]string, mailbox.Error) {
	return b.ListFromPeriodWithContext(context.TODO(), begin, end, limit)
}

// PostWithContext implements mailbox.Box.
func (b *box) PostWithContext(ctx context.Context, id string, c Bytes) (*time.Time, mailbox.Error) {
	b.bf.mu.Lock()
	defer b.bf.mu.Unlock()

	if c == nil {
		return nil, ErrPostingNilContent
	}

	f := b.bf.f

	name := join(f.Name(), id)

	select {
	case <-ctx.Done():
		return nil, mailbox.ErrUnableToPostContent
	case got := <-getFileInfo(name):
		err := got.err
		if err == nil {
			return nil, ErrRepeatedContentIdentifier
		} else if !os.IsNotExist(err) {
			return nil, mailbox.ErrUnableToPostContent
		}
	}

	select {
	case <-ctx.Done():
		os.Remove(name)
		return nil, mailbox.ErrUnableToPostContent
	case got := <-openFile(name):
		if got.err != nil {
			return nil, mailbox.ErrUnableToPostContent
		}
		f = got.data
	}

	ct := time.Now()

	var err mailbox.Error
	select {
	case <-ctx.Done():
		err = mailbox.ErrUnableToPostContent
	case ok := <-writeContent(f, c):
		if !ok {
			err = mailbox.ErrUnableToPostContent
		}
	}
	f.Close()
	if err != nil {
		os.Remove(name)
		return nil, err
	}
	return &ct, nil
}

// Post implements mailbox.Box.
func (b *box) Post(id string, c Bytes) (*time.Time, mailbox.Error) {
	return b.PostWithContext(context.TODO(), id, c)
}
