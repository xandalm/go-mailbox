package external

import (
	"context"
	"time"

	"github.com/xandalm/go-mailbox"
)

type BoxBridge interface {
	HandlePost(ctx context.Context, c mailbox.Bytes) (mailbox.Data, mailbox.Error)
	HandleGet(ctx context.Context, id string) (mailbox.Data, mailbox.Error)
	HandleLazyGet(ctx context.Context, ids ...string) chan mailbox.AttemptData
	HandleListFromPeriod(ctx context.Context, begin, end time.Time, limit int) ([]string, mailbox.Error)
	HandleDelete(ctx context.Context, id string) mailbox.Error
	HandleClean(ctx context.Context) mailbox.Error
}

type box struct {
	b BoxBridge
}

// Clean implements mailbox.Box.
func (b *box) Clean() mailbox.Error {
	return b.b.HandleClean(context.TODO())
}

// CleanWithContext implements mailbox.Box.
func (b *box) CleanWithContext(ctx context.Context) mailbox.Error {
	return b.b.HandleClean(ctx)
}

// Delete implements mailbox.Box.
func (b *box) Delete(id string) mailbox.Error {
	return b.b.HandleDelete(context.TODO(), id)
}

// DeleteWithContext implements mailbox.Box.
func (b *box) DeleteWithContext(ctx context.Context, id string) mailbox.Error {
	return b.b.HandleDelete(ctx, id)
}

// Get implements mailbox.Box.
func (b *box) Get(id string) (mailbox.Data, mailbox.Error) {
	return b.b.HandleGet(context.TODO(), id)
}

// GetWithContext implements mailbox.Box.
func (b *box) GetWithContext(ctx context.Context, id string) (mailbox.Data, mailbox.Error) {
	return b.b.HandleGet(ctx, id)
}

// LazyGet implements mailbox.Box.
func (b *box) LazyGet(ids ...string) chan mailbox.AttemptData {
	return b.b.HandleLazyGet(context.TODO(), ids...)
}

// LazyGetWithContext implements mailbox.Box.
func (b *box) LazyGetWithContext(ctx context.Context, ids ...string) chan mailbox.AttemptData {
	return b.b.HandleLazyGet(ctx, ids...)
}

// ListFromPeriod implements mailbox.Box.
func (b *box) ListFromPeriod(begin, end time.Time, limit int) ([]string, mailbox.Error) {
	return b.b.HandleListFromPeriod(context.TODO(), begin, end, limit)
}

// ListFromPeriodWithContext implements mailbox.Box.
func (b *box) ListFromPeriodWithContext(ctx context.Context, begin, end time.Time, limit int) ([]string, mailbox.Error) {
	return b.b.HandleListFromPeriod(ctx, begin, end, limit)
}

// Post implements mailbox.Box.
func (b *box) Post(c mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	return b.b.HandlePost(context.TODO(), c)
}

// PostWithContext implements mailbox.Box.
func (b *box) PostWithContext(ctx context.Context, c mailbox.Bytes) (mailbox.Data, mailbox.Error) {
	return b.b.HandlePost(ctx, c)
}

func NewBox(b BoxBridge) mailbox.Box {
	return &box{b}
}
