package external

import (
	"context"
	"time"

	mailbox "github.com/xandalm/go-mailbox"
)

type spyBoxBridge struct {
	OnPostCalls           int
	OnGetCalls            int
	OnLazyGetCalls        int
	OnListFromPeriodCalls int
}

// OnPost implements BoxBridge.
func (s *spyBoxBridge) OnPost(ctx context.Context, c mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	s.OnPostCalls++
	return mailbox.Data{}, nil
}

// OnGet implements BoxBridge.
func (s *spyBoxBridge) OnGet(ctx context.Context, id string) (mailbox.Data, mailbox.Error) {
	s.OnGetCalls++
	return mailbox.Data{}, nil
}

func (s *spyBoxBridge) OnLazyGet(ctx context.Context, ids ...string) chan mailbox.AttemptData {
	s.OnLazyGetCalls++
	return nil
}

func (s *spyBoxBridge) OnListFromPeriod(ctx context.Context, begin, end time.Time, limit int) ([]string, mailbox.Error) {
	s.OnListFromPeriodCalls++
	return nil, nil
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
