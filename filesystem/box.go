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
	ErrContentNotFound           = mailbox.NewDetailedError(mailbox.ErrUnableToReadContent, "not found")

	errFileAlreadyExists = errors.New("file already exists")
	errFileNotExist      = errors.New("file not exists")
	errUnableToCheckFile = errors.New("unable to check file")
	errUnableToWriteFile = errors.New("unable to write file")
	errUnableToReadFile  = errors.New("unable to read file")
)

type Bytes = mailbox.Bytes

type rw interface {
	Read(string) ([]byte, error)
	Write(string, []byte) error
}

type rwImpl struct{}

func (s *rwImpl) Read(name string) ([]byte, error) {
	f, err := os.Open(name)
	if os.IsNotExist(err) {
		return nil, errFileNotExist
	}
	if err != nil {
		return nil, errUnableToReadFile
	}
	defer f.Close()
	if data, err := io.ReadAll(f); err == nil {
		return data, nil
	}
	return nil, errUnableToReadFile
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
	if _, err := f.Write(data); err != nil {
		return errUnableToWriteFile
	}
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
	data, err := b.s.Read(filename)
	if err == nil {
		return data, nil
	}
	switch err {
	case errFileNotExist:
		return nil, ErrContentNotFound
	default:
		return nil, mailbox.ErrUnableToReadContent
	}
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
