package filesystem

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
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

type box struct {
	p   *provider
	bf  *boxFile
	idb uint64 // id bias (circular)
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
		if name == idBiasFilename {
			continue
		}
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

func (b *box) newId() string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, b.idb)
	b.idb++

	idbf := b.bf.idbf
	idbf.Seek(0, 0)
	if err := binary.Write(idbf, binary.BigEndian, b.idb); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", sha1.Sum(buf))
}

// PostWithContext implements mailbox.Box.
func (b *box) PostWithContext(ctx context.Context, c Bytes) (mailbox.Data, mailbox.Error) {
	b.bf.mu.Lock()
	defer b.bf.mu.Unlock()

	if c == nil {
		return mailbox.Data{}, ErrPostingNilContent
	}

	f := b.bf.f

	c = bytes.Clone(c)
	ct := time.Now()
	id := b.newId()

	if id == "" {
		return mailbox.Data{}, mailbox.ErrUnableToPostContent
	}

	name := join(f.Name(), id)

	select {
	case <-ctx.Done():
		return mailbox.Data{}, mailbox.ErrUnableToPostContent
	case got := <-getFileInfo(name):
		err := got.err
		if err == nil {
			return mailbox.Data{}, ErrRepeatedContentIdentifier
		} else if !os.IsNotExist(err) {
			return mailbox.Data{}, mailbox.ErrUnableToPostContent
		}
	}

	select {
	case <-ctx.Done():
		os.Remove(name)
		return mailbox.Data{}, mailbox.ErrUnableToPostContent
	case got := <-openFileRW(name):
		if got.err != nil {
			return mailbox.Data{}, mailbox.ErrUnableToPostContent
		}
		f = got.data
	}

	var err mailbox.Error
	select {
	case <-ctx.Done():
		err = mailbox.ErrUnableToPostContent
	case err := <-writeContent(f, c):
		if err != nil {
			err = mailbox.ErrUnableToPostContent
		}
	}
	f.Close()
	if err != nil {
		os.Remove(name)
		return mailbox.Data{}, err
	}
	return mailbox.Data{
		Id:           id,
		CreationTime: ct.UnixNano(),
		Content:      c,
	}, nil
}

// Post implements mailbox.Box.
func (b *box) Post(c Bytes) (mailbox.Data, mailbox.Error) {
	return b.PostWithContext(context.TODO(), c)
}
