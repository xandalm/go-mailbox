package mailbox

import (
	"slices"
	"testing"
)

type stubStorageBox struct {
	id      string
	content map[string]Bytes
}

type stubStorage struct {
	Boxes []stubStorageBox
}

func (s *stubStorage) CreateBox(id string) Error {
	s.Boxes = append(s.Boxes, stubStorageBox{id, make(map[string]Bytes)})
	return nil
}

func (s *stubStorage) ListBoxes() ([]string, Error) {
	ids := []string{}
	for _, b := range s.Boxes {
		ids = append(ids, b.id)
	}
	return ids, nil
}

func (s *stubStorage) DeleteBox(box string) Error {
	s.Boxes = slices.DeleteFunc(s.Boxes, func(b stubStorageBox) bool {
		return b.id == box
	})
	return nil
}

func (s *stubStorage) CleanBox(box string) Error {
	for i := 0; i < len(s.Boxes); i++ {
		b := s.Boxes[i]
		if b.id == box {
			clear(b.content)
		}
	}
	return nil
}

func (s *stubStorage) CreateContent(box string, id string, data []byte) Error {
	for i := 0; i < len(s.Boxes); i++ {
		b := s.Boxes[i]
		if b.id == box {
			if _, ok := b.content[id]; !ok {
				b.content[id] = data
				return nil
			}
			break
		}
	}
	return ErrUnableToPostContent
}

func (s *stubStorage) ReadContent(box string, id string) ([]byte, Error) {
	for i := 0; i < len(s.Boxes); i++ {
		b := s.Boxes[i]
		if b.id == box {
			if data, ok := b.content[id]; ok {
				return data, nil
			}
			break
		}
	}
	return nil, ErrUnableToReadContent
}

func (s *stubStorage) DeleteContent(box string, id string) Error {
	for i := 0; i < len(s.Boxes); i++ {
		b := s.Boxes[i]
		if b.id == box {
			delete(b.content, id)
			return nil
		}
	}
	return ErrUnableToDeleteContent
}

var errFoo Error = newError("foo error")

type stubFailingStorage struct{}

// CleanBox implements Storage.
func (s *stubFailingStorage) CleanBox(string) Error {
	return errFoo
}

// CreateBox implements Storage.
func (s *stubFailingStorage) CreateBox(string) Error {
	return errFoo
}

// CreateContent implements Storage.
func (s *stubFailingStorage) CreateContent(string, string, []byte) Error {
	return errFoo
}

// DeleteBox implements Storage.
func (s *stubFailingStorage) DeleteBox(string) Error {
	return errFoo
}

// DeleteContent implements Storage.
func (s *stubFailingStorage) DeleteContent(string, string) Error {
	return errFoo
}

// List implements Storage.
func (s *stubFailingStorage) ListBoxes() ([]string, Error) {
	return nil, nil
}

// ReadContent implements Storage.
func (s *stubFailingStorage) ReadContent(string, string) ([]byte, Error) {
	return nil, errFoo
}

func AssertContainsFunc[A any, B any](t testing.TB, collec []A, lf B, fn func(e A, lf B) bool) {
	t.Helper()

	for i := 0; i < len(collec); i++ {
		if fn(collec[i], lf) {
			return
		}
	}
	t.Fatalf("there's no %v in %v", lf, collec)
}
