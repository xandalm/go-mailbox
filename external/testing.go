package external

import (
	"context"

	mailbox "github.com/xandalm/go-mailbox"
)

type spyBoxBridge struct {
	OnPostCalls int
}

// OnPost implements BoxBridge.
func (s *spyBoxBridge) OnPost(ctx context.Context, c mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	s.OnPostCalls++
	return mailbox.Data{}, nil
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
