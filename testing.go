package mailbox

import (
	"slices"
	"testing"
)

type stubBox struct {
	Id string
}

func (s *stubBox) Post(any) (any, Error) {
	panic("unimplemented")
}

func (s *stubBox) Get(any) (any, Error) {
	panic("unimplemented")
}

func (s *stubBox) Delete(any) Error {
	panic("unimplemented")
}

func (s *stubBox) Clean() Error {
	panic("unimplemented")
}

type stubProvider struct {
	Boxes []*stubBox
}

func (s *stubProvider) Create(id string) (Box, Error) {
	return &stubBox{id}, nil
}

func (s *stubProvider) Get(id string) (Box, Error) {
	return &stubBox{id}, nil
}

func (s *stubProvider) Delete(id string) Error {
	s.Boxes = slices.DeleteFunc(s.Boxes, func(sb *stubBox) bool {
		return sb.Id == id
	})
	return nil
}

func (s *stubProvider) List() ([]string, Error) {
	return nil, nil
}

var errFoo Error = newError("foo error")

type stubFailingProvider struct{}

func (s *stubFailingProvider) Create(id string) (Box, Error) {
	return nil, errFoo
}

func (s *stubFailingProvider) Get(id string) (Box, Error) {
	return nil, errFoo
}

func (s *stubFailingProvider) Delete(id string) Error {
	return errFoo
}

func (s *stubFailingProvider) List() ([]string, Error) {
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
