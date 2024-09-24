package mailbox

import (
	"os"
	"path/filepath"
)

type fileSystemStorage struct {
	path string
}

func (s *fileSystemStorage) CreateBox(id string) Error {
	path := filepath.Join(s.path, id)
	if _, err := os.Stat(path); err == nil {
		return ErrRepeatedBoxIdentifier
	}
	os.Mkdir(path, 0666)
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
