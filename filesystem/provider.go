package filesystem

import (
	"os"
	"path/filepath"

	"github.com/xandalm/go-mailbox"
)

type box struct {
}

// Clean implements mailbox.Box.
func (b *box) Clean() mailbox.Error {
	panic("unimplemented")
}

// Delete implements mailbox.Box.
func (b *box) Delete(any) mailbox.Error {
	panic("unimplemented")
}

// Get implements mailbox.Box.
func (b *box) Get(any) (any, mailbox.Error) {
	panic("unimplemented")
}

// Post implements mailbox.Box.
func (b *box) Post(any) (any, mailbox.Error) {
	panic("unimplemented")
}

var (
	ErrEmptyBoxIdentifier    = mailbox.NewDetailedError(mailbox.ErrUnableToCreateBox, "identifier can't be empty")
	ErrRepeatedBoxIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToCreateBox, "repeated identifier")
)

type provider struct {
	path string
}

func NewProvider(path, dir string) *provider {
	path = filepath.Join(path, dir)
	err := os.MkdirAll(filepath.Join(path), 0666)
	if err != nil && !os.IsExist(err) {
		panic("unable to create provider")
	}
	p := &provider{path}
	return p
}

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	if id == "" {
		return nil, ErrEmptyBoxIdentifier
	}
	path := filepath.Join(p.path, id)
	if _, err := os.Stat(path); err == nil {
		return nil, ErrRepeatedBoxIdentifier
	} else if !os.IsNotExist(err) {
		return nil, mailbox.ErrUnableToCreateBox
	}
	if err := os.Mkdir(path, 0666); err != nil {
		return nil, mailbox.ErrUnableToCreateBox
	}
	b := &box{}
	return b, nil
}
