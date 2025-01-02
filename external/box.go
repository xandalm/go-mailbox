package external

import (
	"context"
	"time"

	"github.com/xandalm/go-mailbox"
)

type BoxBridge interface {
	OnPost(ctx context.Context, c mailbox.Bytes) (mailbox.Data, mailbox.Error)
	OnGet(ctx context.Context, id string) (mailbox.Data, mailbox.Error)
	OnLazyGet(ctx context.Context, ids ...string) chan mailbox.AttemptData
}

type box struct {
	b BoxBridge
}

// Clean implements mailbox.Box.
func (b *box) Clean() mailbox.Error {
	panic("unimplemented")
}

// CleanWithContext implements mailbox.Box.
func (b *box) CleanWithContext(context.Context) mailbox.Error {
	panic("unimplemented")
}

// Delete implements mailbox.Box.
func (b *box) Delete(string) mailbox.Error {
	panic("unimplemented")
}

// DeleteWithContext implements mailbox.Box.
func (b *box) DeleteWithContext(context.Context, string) mailbox.Error {
	panic("unimplemented")
}

// Get implements mailbox.Box.
func (b *box) Get(id string) (mailbox.Data, mailbox.Error) {
	return b.b.OnGet(context.TODO(), id)
}

// GetWithContext implements mailbox.Box.
func (b *box) GetWithContext(ctx context.Context, id string) (mailbox.Data, mailbox.Error) {
	return b.b.OnGet(ctx, id)
}

// LazyGet implements mailbox.Box.
func (b *box) LazyGet(ids ...string) chan mailbox.AttemptData {
	return b.b.OnLazyGet(context.TODO(), ids...)
}

// LazyGetWithContext implements mailbox.Box.
func (b *box) LazyGetWithContext(ctx context.Context, ids ...string) chan mailbox.AttemptData {
	return b.b.OnLazyGet(ctx, ids...)
}

// ListFromPeriod implements mailbox.Box.
func (b *box) ListFromPeriod(time.Time, time.Time, int) ([]string, mailbox.Error) {
	panic("unimplemented")
}

// ListFromPeriodWithContext implements mailbox.Box.
func (b *box) ListFromPeriodWithContext(context.Context, time.Time, time.Time, int) ([]string, mailbox.Error) {
	panic("unimplemented")
}

// Post implements mailbox.Box.
func (b *box) Post(c mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	return b.b.OnPost(context.TODO(), c)
}

// PostWithContext implements mailbox.Box.
func (b *box) PostWithContext(ctx context.Context, c mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	return b.b.OnPost(ctx, c)
}

func NewBox(b BoxBridge) mailbox.Box {
	return &box{b}
}
