package filesystem

import (
	"errors"
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
		f, err := os.Open(join(path, id))
		if err != nil {
			panic("unable to load existing boxes")
		}
		if err = p.insertBox(&boxFile{id: id, f: f}); err != nil {
			panic("unable to load existing boxes")
		}
	}
	return p
}

func (p *provider) boxPosition(id string) (int, bool) {
	return slices.BinarySearchFunc(p.boxes, id, func(b *boxFile, id string) int {
		return strings.Compare(b.id, id)
	})
}

func (p *provider) createBox(id string) *boxFile {
	p.mu.Lock()
	defer p.mu.Unlock()

	pos, has := p.boxPosition(id)
	if has {
		return nil
	}
	b := &boxFile{
		id: id,
	}
	p.boxes = slices.Insert(p.boxes, pos, b)
	return b
}

func (p *provider) insertBox(b *boxFile) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	pos, has := p.boxPosition(b.id)
	if has {
		return errors.New("identifier already exists")
	}
	p.boxes = slices.Insert(p.boxes, pos, b)
	return nil
}

func (p *provider) removeBox(b *boxFile) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pos, has := p.boxPosition(b.id)
	if !has {
		return
	}
	p.boxes = slices.Delete(p.boxes, pos, pos+1)
}

func (p *provider) getBox(id string) *boxFile {
	p.mu.Lock()
	defer p.mu.Unlock()

	pos, has := p.boxPosition(id)
	if !has {
		return nil
	}
	return p.boxes[pos]
}

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {

	if id == "" {
		return nil, ErrEmptyBoxIdentifier
	}
	bf := p.createBox(id)
	if bf == nil {
		return nil, ErrRepeatedBoxIdentifier
	}
	path := join(p.path, id)
	err := os.Mkdir(path, 0666)
	if err != nil {
		p.removeBox(bf)
		return nil, mailbox.ErrUnableToCreateBox
	}
	f, err := os.Open(path)
	if err != nil {
		p.removeBox(bf)
		return nil, mailbox.ErrUnableToCreateBox
	}
	bf.f = f
	return &box{
		p:  p,
		bf: bf,
	}, nil
}

func (p *provider) Get(id string) (mailbox.Box, mailbox.Error) {

	if bf := p.getBox(id); bf != nil {
		return &box{
			p:  p,
			bf: bf,
		}, nil
	}
	return nil, ErrBoxNotFound
}

func (p *provider) Contains(id string) (bool, mailbox.Error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, has := p.boxPosition(id)
	return has, nil
}

func (p *provider) Delete(id string) mailbox.Error {

	bf := p.getBox(id)

	bf.mu.Lock()
	defer bf.mu.Unlock()

	p.removeBox(bf)

	if err := bf.f.Close(); err != nil {
		p.insertBox(bf)
		return mailbox.ErrUnableToDeleteBox
	}
	if err := os.RemoveAll(join(p.path, id)); err != nil {
		p.insertBox(bf)
		return mailbox.ErrUnableToDeleteBox
	}
	return nil
}
