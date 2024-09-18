package filesystem

import (
	"os"

	"github.com/xandalm/go-mailbox"
)

type box struct {
	p  *provider
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
func (b *box) Post(c any) (any, mailbox.Error) {
	id := "id"
	f, err := os.OpenFile(join(b.p.path, b.id, id), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, mailbox.ErrUnableToPostContent
	}
	defer f.Close()
	return id, nil
}
