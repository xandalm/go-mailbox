package external

import (
	"context"

	"github.com/xandalm/go-mailbox"
)

type ProviderBridge interface {
	OnCreate(context.Context, string) (mailbox.Box, mailbox.Error)
	OnGet(context.Context, string) (mailbox.Box, mailbox.Error)
	OnContains(context.Context, string) bool
	OnDelete(context.Context, string) mailbox.Error
}

type provider struct {
	p ProviderBridge
}

func NewProvider(p ProviderBridge) mailbox.Provider {
	return &provider{p}
}

// Contains implements mailbox.Provider.
func (p *provider) Contains(string) bool {
	panic("unimplemented")
}

// Create implements mailbox.Provider.
func (p *provider) Create(id string) (mailbox.Box, mailbox.Error) {
	return p.p.OnCreate(context.TODO(), id)
}

// Delete implements mailbox.Provider.
func (p *provider) Delete(string) mailbox.Error {
	panic("unimplemented")
}

// Get implements mailbox.Provider.
func (p *provider) Get(string) (mailbox.Box, mailbox.Error) {
	panic("unimplemented")
}
