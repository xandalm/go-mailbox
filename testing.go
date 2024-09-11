package mailbox

import (
	"errors"
)

type stubBox struct {
	Id string
}

type stubProvider struct{}

func (s *stubProvider) Create(id string) (Box, error) {
	return &stubBox{id}, nil
}

func (s *stubProvider) Get(id string) (Box, error) {
	return &stubBox{id}, nil
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

func (s *stubFailingProvider) List() ([]string, error) {
	return nil, errFoo
}
