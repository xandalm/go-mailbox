package external

import (
	"context"
	"time"

	mailbox "github.com/xandalm/go-mailbox"
)

type spyBoxBridge struct {
	HandlePostCalls           int
	HandleGetCalls            int
	HandleLazyGetCalls        int
	HandleListFromPeriodCalls int
	HandleDeleteCalls         int
	HandleCleanCalls          int
}

// HandleClean implements BoxBridge.
func (s *spyBoxBridge) HandleClean(ctx context.Context) mailbox.Error {
	s.HandleCleanCalls++
	return nil
}

// HandlePost implements BoxBridge.
func (s *spyBoxBridge) HandlePost(ctx context.Context, c mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	s.HandlePostCalls++
	return mailbox.Data{}, nil
}

// HandleGet implements BoxBridge.
func (s *spyBoxBridge) HandleGet(ctx context.Context, id string) (mailbox.Data, mailbox.Error) {
	s.HandleGetCalls++
	return mailbox.Data{}, nil
}

// HandleDelete implements BoxBridge.
func (s *spyBoxBridge) HandleLazyGet(ctx context.Context, ids ...string) chan mailbox.AttemptData {
	s.HandleLazyGetCalls++
	return nil
}

// HandleDelete implements BoxBridge.
func (s *spyBoxBridge) HandleListFromPeriod(ctx context.Context, begin, end time.Time, limit int) ([]string, mailbox.Error) {
	s.HandleListFromPeriodCalls++
	return nil, nil
}

// HandleDelete implements BoxBridge.
func (s *spyBoxBridge) HandleDelete(ctx context.Context, id string) mailbox.Error {
	s.HandleDeleteCalls++
	return nil
}

type spyProviderBridge struct {
	HandleCreateCalls   int
	HandleContainsCalls int
	HandleGetCalls      int
	HandleDeleteCalls   int
}

// HandleDelete implements ProviderBridge.
func (s *spyProviderBridge) HandleCreate(context.Context, string) (mailbox.Box, mailbox.Error) {
	s.HandleCreateCalls++
	return nil, nil
}

// HandleDelete implements ProviderBridge.
func (s *spyProviderBridge) HandleDelete(context.Context, string) mailbox.Error {
	s.HandleDeleteCalls++
	return nil
}

// HandleContains implements ProviderBridge.
func (s *spyProviderBridge) HandleContains(context.Context, string) bool {
	s.HandleContainsCalls++
	return false
}

// HandleGet implements ProviderBridge.
func (s *spyProviderBridge) HandleGet(context.Context, string) (mailbox.Box, mailbox.Error) {
	s.HandleGetCalls++
	return nil, nil
}
