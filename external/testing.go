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

type stubProviderBridge struct {
}

// OnDelete implements Provider.
func (s *stubProviderBridge) OnDelete(context.Context, string) mailbox.Error {
	panic("unimplemented")
}

// OnContains implements Provider.
func (s *stubProviderBridge) OnContains(context.Context, string) bool {
	panic("unimplemented")
}

// OnCreate implements Provider.
func (s *stubProviderBridge) OnCreate(context.Context, string) (mailbox.Box, mailbox.Error) {
	return &stubBox{}, nil
}

// OnGet implements Provider.
func (s *stubProviderBridge) OnGet(context.Context, string) (mailbox.Box, mailbox.Error) {
	panic("unimplemented")
}
