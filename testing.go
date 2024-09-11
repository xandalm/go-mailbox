package mailbox

import (
	"errors"
	"slices"
)

type stubBox struct {
	Id string
}

func (s *stubBox) Post(any, any) error {
	panic("unimplemented")
}

func (s *stubBox) Get(any) (any, error) {
	panic("unimplemented")
}

func (s *stubBox) Delete(any) error {
	panic("unimplemented")
}

func (s *stubBox) Clean() error {
	panic("unimplemented")
}

type stubProvider struct {
	Boxes []*stubBox
}

func (s *stubProvider) Create(id string) (Box, error) {
	return &stubBox{id}, nil
}

func (s *stubProvider) Get(id string) (Box, error) {
	return &stubBox{id}, nil
}

func (s *stubProvider) Delete(id string) error {
	s.Boxes = slices.DeleteFunc(s.Boxes, func(sb *stubBox) bool {
		return sb.Id == id
	})
	return nil
}

func (s *stubProvider) List() ([]string, error) {
	return nil, nil
}

var errFoo = errors.New("foo error")

type stubFailingProvider struct{}

func (s *stubFailingProvider) Create(id string) (Box, error) {
	return nil, errFoo
}

func (s *stubFailingProvider) Get(id string) (Box, error) {
	return nil, errFoo
}

func (s *stubFailingProvider) Delete(id string) error {
	return errFoo
}

func (s *stubFailingProvider) List() ([]string, error) {
	return nil, errFoo
}
