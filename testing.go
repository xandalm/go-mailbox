package mailbox

import (
	"context"
	"slices"
	"testing"
	"time"
)

type stubBox struct {
	Id string
}

// CleanWithContext implements Box.
func (s *stubBox) CleanWithContext(context.Context) Error {
	panic("unimplemented")
}

// DeleteWithContext implements Box.
func (s *stubBox) DeleteWithContext(context.Context, string) Error {
	panic("unimplemented")
}

// GetWithContext implements Box.
func (s *stubBox) GetWithContext(context.Context, string) (Data, Error) {
	panic("unimplemented")
}

// LazyGetWithContext implements Box.
func (s *stubBox) LazyGetWithContext(context.Context, ...string) chan AttemptData {
	panic("unimplemented")
}

// ListFromPeriodWithContext implements Box.
func (s *stubBox) ListFromPeriodWithContext(context.Context, time.Time, time.Time, int) ([]string, Error) {
	panic("unimplemented")
}

// PostWithContext implements Box.
func (s *stubBox) PostWithContext(context.Context, string, Bytes) (*time.Time, Error) {
	panic("unimplemented")
}

func (s *stubBox) Post(string, Bytes) (*time.Time, Error) {
	panic("unimplemented")
}

func (s *stubBox) Get(string) (Data, Error) {
	panic("unimplemented")
}

func (s *stubBox) LazyGet(...string) chan AttemptData {
	panic("unimplemented")
}

func (s *stubBox) ListFromPeriod(begin, end time.Time, limit int) ([]string, Error) {
	panic("unimplemented")
}

func (s *stubBox) Delete(string) Error {
	panic("unimplemented")
}

func (s *stubBox) Clean() Error {
	panic("unimplemented")
}

type stubProvider struct {
	Boxes []*stubBox
}

func (s *stubProvider) Create(id string) (Box, Error) {
	b := &stubBox{id}
	s.Boxes = append(s.Boxes, b)
	return b, nil
}

func (s *stubProvider) Get(id string) (Box, Error) {
	return &stubBox{id}, nil
}

func (s *stubProvider) Contains(id string) (bool, Error) {
	has := slices.ContainsFunc(s.Boxes, func(b *stubBox) bool {
		return b.Id == id
	})
	return has, nil
}

func (s *stubProvider) Delete(id string) Error {
	s.Boxes = slices.DeleteFunc(s.Boxes, func(sb *stubBox) bool {
		return sb.Id == id
	})
	return nil
}

func (s *stubProvider) List() ([]string, Error) {
	ret := []string{}
	for i := 0; i < len(s.Boxes); i++ {
		ret = append(ret, s.Boxes[i].Id)
	}
	return ret, nil
}

var errFoo Error = newError("foo error")

type stubFailingProvider struct{}

func (s *stubFailingProvider) Create(id string) (Box, Error) {
	return nil, errFoo
}

func (s *stubFailingProvider) Get(id string) (Box, Error) {
	return nil, errFoo
}

func (s *stubFailingProvider) Contains(id string) (bool, Error) {
	return false, errFoo
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
