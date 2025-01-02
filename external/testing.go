package external

import (
	"context"
	"time"

	mailbox "github.com/xandalm/go-mailbox"
)

type stubBox struct {
}

// Clean implements mailbox.Box.
func (s *stubBox) Clean() mailbox.Error {
	panic("unimplemented")
}

// CleanWithContext implements mailbox.Box.
func (s *stubBox) CleanWithContext(context.Context) mailbox.Error {
	panic("unimplemented")
}

// Delete implements mailbox.Box.
func (s *stubBox) Delete(string) mailbox.Error {
	panic("unimplemented")
}

// DeleteWithContext implements mailbox.Box.
func (s *stubBox) DeleteWithContext(context.Context, string) mailbox.Error {
	panic("unimplemented")
}

// Get implements mailbox.Box.
func (s *stubBox) Get(string) (mailbox.Data, mailbox.Error) {
	panic("unimplemented")
}

// GetWithContext implements mailbox.Box.
func (s *stubBox) GetWithContext(context.Context, string) (mailbox.Data, mailbox.Error) {
	panic("unimplemented")
}

// LazyGet implements mailbox.Box.
func (s *stubBox) LazyGet(...string) chan mailbox.AttemptData {
	panic("unimplemented")
}

// LazyGetWithContext implements mailbox.Box.
func (s *stubBox) LazyGetWithContext(context.Context, ...string) chan mailbox.AttemptData {
	panic("unimplemented")
}

// ListFromPeriod implements mailbox.Box.
func (s *stubBox) ListFromPeriod(time.Time, time.Time, int) ([]string, mailbox.Error) {
	panic("unimplemented")
}

// ListFromPeriodWithContext implements mailbox.Box.
func (s *stubBox) ListFromPeriodWithContext(context.Context, time.Time, time.Time, int) ([]string, mailbox.Error) {
	panic("unimplemented")
}

// Post implements mailbox.Box.
func (s *stubBox) Post(mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	panic("unimplemented")
}

// PostWithContext implements mailbox.Box.
func (s *stubBox) PostWithContext(context.Context, mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	panic("unimplemented")
}

type spyProviderBridge struct {
	OnCreateCalls   int
	OnContainsCalls int
	OnGetCalls      int
	OnDeleteCalls   int
}

// OnDelete implements ProviderBridge.
func (s *spyProviderBridge) OnCreate(context.Context, string) (mailbox.Box, mailbox.Error) {
	s.OnCreateCalls++
	return nil, nil
}

// OnDelete implements ProviderBridge.
func (s *spyProviderBridge) OnDelete(context.Context, string) mailbox.Error {
	s.OnDeleteCalls++
	return nil
}

// OnContains implements ProviderBridge.
func (s *spyProviderBridge) OnContains(context.Context, string) bool {
	s.OnContainsCalls++
	return false
}

// OnGet implements ProviderBridge.
func (s *spyProviderBridge) OnGet(context.Context, string) (mailbox.Box, mailbox.Error) {
	s.OnGetCalls++
	return nil, nil
}
