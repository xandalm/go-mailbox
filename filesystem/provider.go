package filesystem

import (
	"os"
	"path/filepath"

	"github.com/xandalm/go-mailbox"
)

type box struct {
	id string
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
	ErrBoxNotFound           = mailbox.NewDetailedError(mailbox.ErrUnableToRestoreBox, "not found")
)

func join(s ...string) string {
	return filepath.Join(s...)
}

type provider struct {
	path string
}

func NewProvider(path, dir string) *provider {
	path = join(path, dir)
	err := os.MkdirAll(path, 0666)
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
	path := join(p.path, id)
	if _, err := os.Stat(path); err == nil {
		return nil, ErrRepeatedBoxIdentifier
	} else if !os.IsNotExist(err) {
		return nil, mailbox.ErrUnableToCreateBox
	}
	if err := os.Mkdir(path, 0666); err != nil {
		return nil, mailbox.ErrUnableToCreateBox
	}
	b := &box{id}
	return b, nil
}

func (p *provider) Get(id string) (mailbox.Box, mailbox.Error) {
	if _, err := os.Stat(join(p.path, id)); err == nil {
		return &box{id}, nil
	} else if os.IsNotExist(err) {
		return nil, ErrBoxNotFound
	}
	return nil, mailbox.ErrUnableToRestoreBox
}

func (p *provider) Delete(id string) mailbox.Error {
	if err := os.RemoveAll(join(p.path, id)); err != nil {
		return mailbox.ErrUnableToDeleteBox
	}
	return nil
}

func (p *provider) List() ([]string, mailbox.Error) {
	de, err := os.ReadDir(p.path)
	if err != nil {
		return nil, mailbox.ErrUnableToRestoreBox
	}
	ret := []string{}
	for i := 0; i < len(de); i++ {
		ret = append(ret, de[i].Name())
	}
	return ret, nil
}
