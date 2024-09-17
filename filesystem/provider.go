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

type provider struct {
	path string
}

func NewProvider(path, dir string) *provider {
	path = filepath.Join(path, dir)
	p := &provider{path}
	err := os.MkdirAll(filepath.Join(p.path), 0666)
	if err != nil && !os.IsExist(err) {
		panic("unable to create provider")
	}
	return p
}

func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	b := &box{}
	os.Mkdir(filepath.Join(p.path, id), 0666)
	return b, nil
}
