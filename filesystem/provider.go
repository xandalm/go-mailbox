package filesystem

import (
	"os"
	"path/filepath"
	"slices"

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
	foundBoxes, _ := f.Readdirnames(0)
	for _, id := range foundBoxes {
		p.boxes = append(p.boxes, &box{&fsHandlerImpl{}, p, id})
	}
	return p
}

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	if id == "" {
		return nil, ErrEmptyBoxIdentifier
	}
	path := join(p.path, id)
	if _, err := os.Stat(path); err == nil {
		return nil, ErrRepeatedBoxIdentifier
	} else if !os.IsNotExist(err) {
		return nil, mailbox.ErrUnableToCreateBox
	}
	if err := os.Mkdir(path, 0666); err != nil {
		return nil, mailbox.ErrUnableToCreateBox
	}
	b := &box{&fsHandlerImpl{}, p, id}
	return b, nil
}

func (p *provider) Get(id string) (mailbox.Box, mailbox.Error) {
	if _, err := os.Stat(join(p.path, id)); err == nil {
		return &box{&fsHandlerImpl{}, p, id}, nil
	} else if os.IsNotExist(err) {
		return nil, ErrBoxNotFound
	}
	return nil, mailbox.ErrUnableToRestoreBox
}

func (p *provider) Contains(id string) (bool, mailbox.Error) {
	return slices.ContainsFunc(p.boxes, func(b *box) bool {
		return b.id == id
	}), nil
}

func (p *provider) Delete(id string) mailbox.Error {
	if err := os.RemoveAll(join(p.path, id)); err != nil {
		return mailbox.ErrUnableToDeleteBox
	}
	return nil
}
