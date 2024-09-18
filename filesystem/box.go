package filesystem

import (
	"os"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrRepeatedContentIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "provided identifier is already in use")
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
func (b *box) Delete(string) mailbox.Error {
	panic("unimplemented")
}

// Get implements mailbox.Box.
func (b *box) Get(string) (any, mailbox.Error) {
	panic("unimplemented")
}

// Post implements mailbox.Box.
func (b *box) Post(id string, c any) mailbox.Error {
	filename := join(b.p.path, b.id, id)
	if _, err := os.Stat(filename); err == nil {
		return ErrRepeatedContentIdentifier
	} else if !os.IsNotExist(err) {
		return mailbox.ErrUnableToPostContent
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return mailbox.ErrUnableToPostContent
	}
	defer f.Close()
	return nil
}
