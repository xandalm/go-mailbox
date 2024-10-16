package filesystem

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrEmptyBoxIdentifier    = mailbox.NewDetailedError(mailbox.ErrUnableToCreateBox, "identifier can't be empty")
	ErrRepeatedBoxIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToCreateBox, "repeated identifier")
	ErrBoxNotFound           = mailbox.NewDetailedError(mailbox.ErrUnableToRestoreBox, "not found")
)

func join(s ...string) string {
	return filepath.Join(s...)
}

type boxFile struct {
	mu sync.RWMutex
	id string
	f  *os.File
}

type provider struct {
	mu    sync.RWMutex
	f     *os.File
	boxes []*boxFile
	path  string
}

func NewProvider(path, dir string) mailbox.Provider {
	path = join(path, dir)
	err := os.MkdirAll(path, 0666)
	if err != nil && !os.IsExist(err) {
		panic("unable to create provider")
	}
	f, err := os.Open(path)
	if err != nil {
		panic("unable to keep directory file open")
	}
	p := &provider{sync.RWMutex{}, f, []*boxFile{}, path}
	foundBoxes, err := f.Readdirnames(0)
	if err != nil {
		panic("unable to load existing boxes")
	}
	for _, id := range foundBoxes {
		pos, _ := p.boxPosition(id)
		f, err := os.Open(join(path, id))
		if err != nil {
			panic("unable to load existing boxes")
		}
		p.insertBoxAt(pos, &boxFile{
			id: id,
			f:  f,
		})
	}
	return p
}

func (p *provider) boxPosition(id string) (int, bool) {
	return slices.BinarySearchFunc(p.boxes, id, func(b *boxFile, id string) int {
		return strings.Compare(b.id, id)
	})
}

func (p *provider) insertBoxAt(pos int, b *boxFile) {
	p.boxes = slices.Insert(p.boxes, pos, b)
}

func (p *provider) removeBoxAt(pos int) {
	p.boxes = slices.Delete(p.boxes, pos, pos+1)
}

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if id == "" {
		return nil, ErrEmptyBoxIdentifier
	}
	pos, has := p.boxPosition(id)
	if has {
		return nil, ErrRepeatedBoxIdentifier
	}
	path := join(p.path, id)
	if err := os.Mkdir(path, 0666); err != nil {
		return nil, mailbox.ErrUnableToCreateBox
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, mailbox.ErrUnableToCreateBox
	}
	bf := &boxFile{
		id: id,
		f:  f,
	}
	p.insertBoxAt(pos, bf)
	return &box{
		p:  p,
		bf: bf,
	}, nil
}

func (p *provider) Get(id string) (mailbox.Box, mailbox.Error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	pos, has := p.boxPosition(id)
	if !has {
		return nil, ErrBoxNotFound
	}
	bf := p.boxes[pos]
	return &box{
		p:  p,
		bf: bf,
	}, nil
}

func (p *provider) Contains(id string) (bool, mailbox.Error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, has := p.boxPosition(id)
	return has, nil
}

func (p *provider) Delete(id string) mailbox.Error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pos, has := p.boxPosition(id)
	if !has {
		return nil
	}
	bf := p.boxes[pos]

	bf.mu.Lock()
	defer bf.mu.Unlock()

	if err := bf.f.Close(); err != nil {
		return mailbox.ErrUnableToDeleteBox
	}
	if err := os.RemoveAll(join(p.path, id)); err != nil {
		return mailbox.ErrUnableToDeleteBox
	}
	p.removeBoxAt(pos)
	return nil
}
