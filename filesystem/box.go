package filesystem

import (
	"os"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrRepeatedContentIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "provided identifier is already in use")
	ErrPostingNilContent         = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "can't post nil content")
	ErrContentNotFound           = mailbox.NewDetailedError(mailbox.ErrUnableToReadContent, "not found")
)

type Bytes = mailbox.Bytes

type box struct {
	f  *os.File
	p  *provider
	id string
}

// Clean implements mailbox.Box.
func (b *box) Clean() mailbox.Error {
	names, err := b.f.Readdirnames(0)
	if err != nil {
		return mailbox.ErrUnableToCleanBox
	}
	for _, name := range names {
		os.Remove(join(b.f.Name(), name))
	}
	return nil
}

// Delete implements mailbox.Box.
func (b *box) Delete(id string) mailbox.Error {
	name := join(b.f.Name(), id)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return mailbox.ErrUnableToReadContent
	}
	if err := os.Remove(name); err != nil {
		return mailbox.ErrUnableToDeleteContent
	}
	return nil
}

// Get implements mailbox.Box.
func (b *box) Get(id string) (Bytes, mailbox.Error) {
	name := join(b.f.Name(), id)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return nil, ErrContentNotFound
	} else if err != nil {
		return nil, mailbox.ErrUnableToReadContent
	}
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, mailbox.ErrUnableToReadContent
	}
	return data, nil
}

// Post implements mailbox.Box.
func (b *box) Post(id string, c Bytes) mailbox.Error {
	if c == nil {
		return ErrPostingNilContent
	}
	name := join(b.f.Name(), id)
	if _, err := os.Stat(name); err == nil {
		return ErrRepeatedContentIdentifier
	} else if !os.IsNotExist(err) {
		return mailbox.ErrUnableToPostContent
	}
	err := os.WriteFile(name, c, 0666)
	if err != nil {
		return mailbox.ErrUnableToPostContent
	}
	return nil
}
