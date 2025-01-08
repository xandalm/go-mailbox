package external

import (
	"context"

	"github.com/xandalm/go-mailbox"
)

type ProviderBridge interface {
	HandleCreate(context.Context, string) (mailbox.Box, mailbox.Error)
	HandleGet(context.Context, string) (mailbox.Box, mailbox.Error)
	HandleContains(context.Context, string) bool
	HandleDelete(context.Context, string) mailbox.Error
}

type provider struct {
	p ProviderBridge
}

func NewProvider(p ProviderBridge) mailbox.Provider {
	return &provider{p}
}

// Contains implements mailbox.Provider.
func (p *provider) Contains(ctx context.Context, id string) bool {
	return p.p.HandleContains(ctx, id)
}

// Create implements mailbox.Provider.
func (p *provider) Create(ctx context.Context, id string) (mailbox.Box, mailbox.Error) {
	return p.p.HandleCreate(ctx, id)
}

// Delete implements mailbox.Provider.
func (p *provider) Delete(ctx context.Context, id string) mailbox.Error {
	return p.p.HandleDelete(ctx, id)
}

// Get implements mailbox.Provider.
func (p *provider) Get(ctx context.Context, id string) (mailbox.Box, mailbox.Error) {
	return p.p.HandleGet(ctx, id)
}
