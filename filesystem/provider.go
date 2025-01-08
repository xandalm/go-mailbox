package filesystem

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
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
)

const idBiasFilename = ".idb"

func join(s ...string) string {
	return filepath.Join(s...)
}

type boxFile struct {
	mu   sync.RWMutex
	id   string
	f    *os.File
	idbf *os.File // id bias file
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
		panic(fmt.Sprintf("unable to create provider, %v", err))
	}
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("unable to keep directory file open, %v", err))
	}
	p := &provider{sync.RWMutex{}, f, []*boxFile{}, path}
	foundBoxes, err := f.Readdirnames(0)
	if err != nil {
		panic(fmt.Sprintf("unable to load existing boxes, %v", err))
	}
	for _, id := range foundBoxes {
		box := &boxFile{id: id}
		box.f, err = os.Open(join(path, id))
		if err != nil {
			panic(fmt.Sprintf("unable to load existing boxes, %v", err))
		}
		box.idbf, err = os.OpenFile(join(path, idBiasFilename), os.O_RDWR, 0666)
		if err == nil {
			box.f.Close()
			panic(fmt.Sprintf("unable to load existing boxes, %v", err))
		}
		if err = p.insertBox(box); err != nil {
			box.f.Close()
			box.idbf.Close()
			panic(fmt.Sprintf("unable to load existing boxes, %v", err))
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

func (p *provider) Create(ctx context.Context, id string) (mailbox.Box, mailbox.Error) {

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

	var idb uint64
	f, err = os.OpenFile(join(path, idBiasFilename), os.O_CREATE|os.O_RDWR, 0666)
	if err == nil {
		err = binary.Write(f, binary.BigEndian, idb)
	}
	if err != nil {
		bf.f.Close()
		os.Remove(path)
		p.removeBox(bf)
		return nil, mailbox.ErrUnableToCreateBox
	}
	bf.idbf = f

	return &box{
		p:   p,
		bf:  bf,
		idb: idb,
	}, nil
}

func (p *provider) Get(ctx context.Context, id string) (mailbox.Box, mailbox.Error) {
	if bf := p.getBox(id); bf != nil {
		return &box{
			p:  p,
			bf: bf,
		}, nil
	}
	return nil, mailbox.ErrBoxNotFound
}

func (p *provider) Contains(ctx context.Context, id string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, has := p.boxPosition(id)
	return has
}

func (p *provider) Delete(ctx context.Context, id string) mailbox.Error {

	bf := p.getBox(id)

	bf.mu.Lock()
	defer bf.mu.Unlock()

	p.removeBox(bf)

	if err := bf.f.Close(); err != nil {
		p.insertBox(bf)
		return mailbox.ErrUnableToDeleteBox
	}
	if err := bf.idbf.Close(); err != nil {
		p.insertBox(bf)
		return mailbox.ErrUnableToDeleteBox
	}
	if err := os.RemoveAll(join(p.path, id)); err != nil {
		p.insertBox(bf)
		return mailbox.ErrUnableToDeleteBox
	}
	return nil
}
