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

type provider struct {
	mu    sync.RWMutex
	f     *os.File
	boxes []*box
	path  string
}

func NewProvider(path, dir string) *provider {
	path = join(path, dir)
	err := os.MkdirAll(path, 0666)
	if err != nil && !os.IsExist(err) {
		panic("unable to create provider")
	}
	f, err := os.Open(path)
	if err != nil {
		panic("unable to keep directory file open")
	}
	p := &provider{sync.RWMutex{}, f, []*box{}, path}
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
		p.insertBoxAt(pos, &box{f, p, id})
	}
	return p
}

func (p *provider) boxPosition(id string) (int, bool) {
	return slices.BinarySearchFunc(p.boxes, id, func(b *box, id string) int {
		return strings.Compare(b.id, id)
	})
}

func (p *provider) insertBoxAt(pos int, b *box) {
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
	b := &box{f, p, id}
	p.insertBoxAt(pos, b)
	return b, nil
}

func (p *provider) Get(id string) (mailbox.Box, mailbox.Error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	pos, has := p.boxPosition(id)
	if !has {
		return nil, ErrBoxNotFound
	}
	return p.boxes[pos], nil
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
	if err := p.boxes[pos].f.Close(); err != nil {
		return mailbox.ErrUnableToDeleteBox
	}
	if err := os.RemoveAll(join(p.path, id)); err != nil {
		return mailbox.ErrUnableToDeleteBox
	}
	p.removeBoxAt(pos)
	return nil
}
