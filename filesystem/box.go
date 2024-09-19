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

	errFileAlreadyExists   = errors.New("file already exists")
	errFileNotExist        = errors.New("file not exists")
	errUnableToCheckFile   = errors.New("unable to check file")
	errUnableToWriteFile   = errors.New("unable to write file")
	errUnableToReadFile    = errors.New("unable to read file")
	errUnableToDeleteFile  = errors.New("unable to delete file")
	errUnableToCleanFolder = errors.New("unable to clean folder files")
)

type Bytes = mailbox.Bytes

type fsHandler interface {
	Read(string, string) ([]byte, error)
	Write(string, string, []byte) error
	Delete(string, string) error
	Clean(string) error
}

type fsHandlerImpl struct{}

func (fs *fsHandlerImpl) Read(path, id string) ([]byte, error) {
	name := join(path, id)
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

func (fs *fsHandlerImpl) Write(path, id string, data []byte) error {
	name := join(path, id)
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

func (fs *fsHandlerImpl) Delete(path, id string) error {
	name := join(path, id)
	if err := os.Remove(name); err != nil && !os.IsNotExist(err) {
		return errUnableToDeleteFile
	}
	return nil
}

func (fs *fsHandlerImpl) Clean(path string) error {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return errUnableToCleanFolder
	}
	for _, e := range dirEntry {
		os.Remove(join(path, e.Name()))
	}
	return nil
}

type box struct {
	fs fsHandler
	p  *provider
	id string
}

func (b *box) path() string {
	return join(b.p.path, b.id)
}

// Clean implements mailbox.Box.
func (b *box) Clean() mailbox.Error {
	if err := b.fs.Clean(b.path()); err != nil {
		return mailbox.ErrUnableToCleanBox
	}
	return nil
}

// Delete implements mailbox.Box.
func (b *box) Delete(id string) mailbox.Error {
	if err := b.fs.Delete(b.path(), id); err != nil {
		return mailbox.ErrUnableToDeleteContent
	}
	return nil
}

// Get implements mailbox.Box.
func (b *box) Get(id string) (Bytes, mailbox.Error) {
	data, err := b.fs.Read(b.path(), id)
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
	err := b.fs.Write(b.path(), id, c)
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
