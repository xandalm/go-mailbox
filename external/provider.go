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
func (p *provider) Contains(id string) bool {
	return p.p.HandleContains(context.TODO(), id)
}

// Create implements mailbox.Provider.
func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	return p.p.HandleCreate(context.TODO(), id)
}

// Delete implements mailbox.Provider.
func (p *provider) Delete(id string) mailbox.Error {
	return p.p.HandleDelete(context.TODO(), id)
}

// Get implements mailbox.Provider.
func (p *provider) Get(id string) (mailbox.Box, mailbox.Error) {
	return p.p.HandleGet(context.TODO(), id)
}
