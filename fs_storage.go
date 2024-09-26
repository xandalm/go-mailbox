package mailbox

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type boxInfo struct {
	id string
	f  *os.File
}

type fileSystemStorage struct {
	mu        sync.RWMutex
	f         *os.File
	boxesInfo []*boxInfo
}

func (s *fileSystemStorage) searchBoxPosition(id string) (int, bool) {
	return slices.BinarySearchFunc(s.boxesInfo, id, func(e *boxInfo, id string) int {
		return strings.Compare(e.id, id)
	})
}

func (s *fileSystemStorage) CreateBox(id string) Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.f.Name(), id)
	pos, has := s.searchBoxPosition(id)
	if has {
		return ErrRepeatedBoxIdentifier
	}
	if err := os.Mkdir(path, 0666); err != nil {
		return ErrUnableToCreateBox
	}
	if f, err := os.Open(path); err == nil {
		s.boxesInfo = slices.Insert(
			s.boxesInfo,
			pos,
			&boxInfo{id, f},
		)
		return nil
	}
	return ErrUnableToCreateBox
}

func (s *fileSystemStorage) ListBoxes() ([]string, Error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := []string{}
	for _, box := range s.boxesInfo {
		ids = append(ids, box.id)
	}
	return ids, nil
}

func (s *fileSystemStorage) DeleteBox(id string) Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos, has := s.searchBoxPosition(id)
	if !has {
		return nil
	}
	box := s.boxesInfo[pos]
	if box.f.Close() != nil || os.RemoveAll(box.f.Name()) != nil {
		return ErrUnableToDeleteBox
	}
	s.boxesInfo = slices.Delete(s.boxesInfo, pos, pos+1)
	return nil
}

func (s *fileSystemStorage) CleanBox(id string) Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos, has := s.searchBoxPosition(id)
	if !has {
		return nil
	}
	f := s.boxesInfo[pos].f
	contents, err := f.Readdirnames(0)
	if err != nil {
		return ErrUnableToCleanBox
	}
	for _, c := range contents {
		err = os.Remove(filepath.Join(f.Name(), c))
		if err == nil {
			continue
		}
		return ErrUnableToCleanBox
	}
	return nil
}

func (s *fileSystemStorage) CreateContent(bid, cid string, d []byte) Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos, has := s.searchBoxPosition(bid)
	if !has {
		return ErrBoxNotFoundToPostContent
	}
	path := filepath.Join(s.boxesInfo[pos].f.Name(), cid)
	if _, err := os.Stat(path); err == nil {
		return ErrRepeatedContentIdentifier
	} else if !os.IsNotExist(err) {
		return ErrUnableToPostContent
	}
	if err := os.WriteFile(path, d, 0666); err != nil {
		return ErrUnableToPostContent
	}
	return nil
}

func (s *fileSystemStorage) ReadContent(bid, cid string) ([]byte, Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos, has := s.searchBoxPosition(bid)
	if !has {
		return nil, ErrBoxNotFoundToReadContent
	}
	path := filepath.Join(s.boxesInfo[pos].f.Name(), cid)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, ErrContentNotFound
	} else if err != nil {
		return nil, ErrUnableToReadContent
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrUnableToReadContent
	}
	return data, nil
}

func (s *fileSystemStorage) DeleteContent(bid, cid string) Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos, has := s.searchBoxPosition(bid)
	if !has {
		return nil
	}
	path := filepath.Join(s.boxesInfo[pos].f.Name(), cid)
	if err := os.Remove(path); err != nil {
		return ErrUnableToDeleteContent
	}
	return nil
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
		f:         f,
		boxesInfo: []*boxInfo{},
	}
}
