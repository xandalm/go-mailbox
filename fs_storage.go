package mailbox

import (
	"container/list"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type fileSystemHandler interface {
	Open(string) (*os.File, error)
	Exists(string) (bool, error)
	Mkdir(string) error
	Ls(string) ([]string, error)
	Remove(string) error
	Clean(string) error
	WriteFile(string, []byte) error
	ReadFile(string) ([]byte, error)
}

type defaulFileSystemHandler struct{}

func (h *defaulFileSystemHandler) Open(name string) (*os.File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (h *defaulFileSystemHandler) Exists(name string) (bool, error) {
	if _, err := os.Stat(name); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (h *defaulFileSystemHandler) Mkdir(name string) error {
	return os.Mkdir(name, 0666)
}

func (h *defaulFileSystemHandler) Ls(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (h *defaulFileSystemHandler) Remove(name string) error {
	return os.Remove(name)
}

func (h *defaulFileSystemHandler) Clean(name string) error {
	names, err := h.Ls(name)
	if err != nil {
		return err
	}
	for _, n := range names {
		err = h.Remove(filepath.Join(name, n))
		if err != nil {
			break
		}
	}
	return err
}

func (h *defaulFileSystemHandler) WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0666)
}

func (h *defaulFileSystemHandler) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

type idxNode struct {
	id string
	e  *list.Element
}

type fileSystemStorage struct {
	f           *os.File
	boxFiles    *list.List
	boxFilesIdx []*idxNode
	handler     fileSystemHandler
	path        string
}

func (s *fileSystemStorage) searchBoxPosition(id string) (int, bool) {
	return slices.BinarySearchFunc(s.boxFilesIdx, id, func(e *idxNode, id string) int {
		return strings.Compare(e.id, id)
	})
}

func (s *fileSystemStorage) CreateBox(id string) Error {
	path := filepath.Join(s.path, id)
	pos, has := s.searchBoxPosition(id)
	if has {
		return ErrRepeatedBoxIdentifier
	}
	if err := s.handler.Mkdir(path); err != nil {
		return ErrUnableToCreateBox
	}
	if f, err := s.handler.Open(path); err == nil {
		e := s.boxFiles.PushBack(f)
		s.boxFilesIdx = slices.Insert(
			s.boxFilesIdx,
			pos,
			&idxNode{id, e},
		)
		return nil
	}
	return ErrUnableToCreateBox
}

func (s *fileSystemStorage) ListBoxes() ([]string, Error) {
	if ids, err := s.handler.Ls(s.path); err != nil {
		return nil, ErrUnableToListBoxes
	} else {
		return ids, nil
	}
}

func (s *fileSystemStorage) DeleteBox(id string) Error {
	path := filepath.Join(s.path, id)
	pos, has := s.searchBoxPosition(id)
	if !has {
		return nil
	}
	node := s.boxFilesIdx[pos]
	bf := node.e.Value.(*os.File)
	if bf.Close() != nil || s.handler.Remove(path) != nil {
		return ErrUnableToDeleteBox
	}
	s.boxFiles.Remove(node.e)
	s.boxFilesIdx = slices.Delete(s.boxFilesIdx, pos, pos+1)
	return nil
}

func (s *fileSystemStorage) CleanBox(id string) Error {
	_, has := s.searchBoxPosition(id)
	if !has {
		return nil
	}
	path := filepath.Join(s.path, id)
	if err := s.handler.Clean(path); err != nil {
		return ErrUnableToCleanBox
	}
	return nil
}

func (s *fileSystemStorage) CreateContent(bid, cid string, d []byte) Error {
	path := filepath.Join(s.path, bid, cid)
	exists, err := s.handler.Exists(path)
	if err != nil {
		return ErrUnableToPostContent
	}
	if exists {
		return ErrRepeatedContentIdentifier
	}
	if err := s.handler.WriteFile(path, d); err != nil {
		return ErrUnableToPostContent
	}
	return nil
}

func (s *fileSystemStorage) ReadContent(bid, cid string) ([]byte, Error) {
	path := filepath.Join(s.path, bid, cid)
	exists, err := s.handler.Exists(path)
	if err != nil {
		return nil, ErrUnableToReadContent
	}
	if !exists {
		return nil, ErrContentNotFound
	}
	data, err := s.handler.ReadFile(path)
	if err != nil {
		return nil, ErrUnableToReadContent
	}
	return data, nil
}

func NewFileSystemStorage(path, dir string) *fileSystemStorage {
	if dir == "" {
		panic("mailbox: loading filesystem storage, the dir (folder name) is required and can't be empty string")
	}

	path = filepath.Join(path, dir)
	err := os.MkdirAll(path, 0666)
	if err != nil && !os.IsExist(err) {
		panic("mailbox: loading filesystem storage, unable to create folder")
	}

	f, err := os.Open(path)
	if err != nil {
		panic("mailbox: loading filesystem storage, unable to open storage file")
	}

	return &fileSystemStorage{
		f:           f,
		boxFiles:    list.New(),
		boxFilesIdx: []*idxNode{},
		handler:     &defaulFileSystemHandler{},
		path:        path,
	}
}
