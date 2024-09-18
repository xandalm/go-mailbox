package filesystem

import (
	"io"
	"os"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrRepeatedContentIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "provided identifier is already in use")
	ErrPostingNilContent         = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "can't post nil content")
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
func (b *box) Get(id string) (any, mailbox.Error) {
	filename := join(b.p.path, b.id, id)
	f, _ := os.Open(filename)
	defer f.Close()
	data, _ := io.ReadAll(f)
	return string(data), nil
}

// Post implements mailbox.Box.
func (b *box) Post(id string, c any) mailbox.Error {
	if c == nil {
		return ErrPostingNilContent
	}
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
