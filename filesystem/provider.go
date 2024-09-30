package filesystem

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

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
	p := &provider{f, []*box{}, path}
	foundBoxes, err := f.Readdirnames(0)
	if err != nil {
		panic("unable to load existing boxes")
	}
	for _, id := range foundBoxes {
		pos, _ := p.boxPosition(id)
		p.insertBoxAt(pos, &box{&fsHandlerImpl{}, p, id})
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

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	if id == "" {
		return nil, ErrEmptyBoxIdentifier
	}
	path := join(p.path, id)
	pos, has := p.boxPosition(id)
	if has {
		return nil, ErrRepeatedBoxIdentifier
	}
	if err := os.Mkdir(path, 0666); err != nil {
		return nil, mailbox.ErrUnableToCreateBox
	}
	b := &box{&fsHandlerImpl{}, p, id}
	p.insertBoxAt(pos, b)
	return b, nil
}

func (p *provider) Get(id string) (mailbox.Box, mailbox.Error) {
	pos, has := p.boxPosition(id)
	if !has {
		return nil, ErrBoxNotFound
	}
	return p.boxes[pos], nil
}

func (p *provider) Contains(id string) (bool, mailbox.Error) {
	_, has := p.boxPosition(id)
	return has, nil
}

func (p *provider) Delete(id string) mailbox.Error {
	if err := os.RemoveAll(join(p.path, id)); err != nil {
		return mailbox.ErrUnableToDeleteBox
	}
	return nil
}
