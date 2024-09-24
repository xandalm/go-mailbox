package mailbox

import (
	"os"
	"path/filepath"
)

type fileSystemHandler interface {
	Exists(string) (bool, error)
	Mkdir(string) error
}

type defaulFileSystemHandler struct{}

func (h *defaulFileSystemHandler) Exists(file string) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (h *defaulFileSystemHandler) Mkdir(dirname string) error {
	return os.Mkdir(dirname, 0666)
}

type fileSystemStorage struct {
	handler fileSystemHandler
	path    string
}

func (s *fileSystemStorage) CreateBox(id string) Error {
	path := filepath.Join(s.path, id)
	if exists, err := s.handler.Exists(path); err != nil {
		return ErrUnableToCreateBox
	} else if exists {
		return ErrRepeatedBoxIdentifier
	}
	if err := s.handler.Mkdir(path); err != nil {
		return ErrUnableToCreateBox
	}
	return nil
}

func NewFileSystemStorage(path, dir string) *fileSystemStorage {
	if dir == "" {
		panic("mailbox: on filesystem storage creation, the dir (folder name) is required and can't be empty string")
	}

	path = filepath.Join(path, dir)
	err := os.MkdirAll(path, 0666)
	if err != nil && !os.IsExist(err) {
		panic("mailbox: on filesystem storage creation, unable to create storage folder")
	}

	return &fileSystemStorage{
		path: path,
	}
}
