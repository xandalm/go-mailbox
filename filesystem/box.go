package filesystem

import (
	"errors"
	"io"
	"os"

	"github.com/xandalm/go-mailbox"
)

var (
	ErrRepeatedContentIdentifier = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "provided identifier is already in use")
	ErrPostingNilContent         = mailbox.NewDetailedError(mailbox.ErrUnableToPostContent, "can't post nil content")

	errFileAlreadyExists = errors.New("file already exists")
	errUnableToCheckFile = errors.New("unable to check file")
	errUnableToWriteFile = errors.New("unable to write file")
)

type Bytes = mailbox.Bytes

type rw interface {
	Read(string) ([]byte, error)
	Write(string, []byte) error
}

type rwImpl struct{}

func (s *rwImpl) Read(name string) ([]byte, error) {
	f, _ := os.Open(name)
	defer f.Close()
	data, _ := io.ReadAll(f)
	return data, nil
}

func (s rwImpl) Write(name string, data []byte) error {
	if _, err := os.Stat(name); err == nil {
		return errFileAlreadyExists
	} else if !os.IsNotExist(err) {
		return errUnableToCheckFile
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return errUnableToWriteFile
	}
	defer f.Close()
	f.Write(data)
	return nil
}

type box struct {
	s  rw
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
func (b *box) Get(id string) (Bytes, mailbox.Error) {
	filename := join(b.p.path, b.id, id)
	data, _ := b.s.Read(filename)
	return data, nil
}

// Post implements mailbox.Box.
func (b *box) Post(id string, c Bytes) mailbox.Error {
	if c == nil {
		return ErrPostingNilContent
	}
	filename := join(b.p.path, b.id, id)
	err := b.s.Write(filename, c)
	if err == nil {
		return nil
	}
	switch err {
	case errFileAlreadyExists:
		return ErrRepeatedContentIdentifier
	default:
		return mailbox.ErrUnableToPostContent
	}
}
