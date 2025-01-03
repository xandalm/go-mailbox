package external

import (
	"context"
	"testing"
	"time"

	"github.com/xandalm/go-testing/assert"
)

func TestNewBox(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	assert.NotNil(t, b)
}

var dummyContent = []byte("abc")

func TestBox_PostHasForwardedToHandlePost(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	t.Run("providing context", func(t *testing.T) {
		b.PostWithContext(context.TODO(), dummyContent)

		assert.Equal(t, bb.HandlePostCalls, 1)
	})
	t.Run("not providing context", func(t *testing.T) {
		b.Post(dummyContent)

		assert.Equal(t, bb.HandlePostCalls, 2)
	})
}

func TestBox_GetHasForwardedToHandleGet(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	t.Run("providing context", func(t *testing.T) {
		b.GetWithContext(context.TODO(), "post")

		assert.Equal(t, bb.HandleGetCalls, 1)
	})
	t.Run("not providing context", func(t *testing.T) {
		b.Get("post")

		assert.Equal(t, bb.HandleGetCalls, 2)
	})
}

func TestBox_LazyGetHasForwardedToHandleLazyGet(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	t.Run("providing context", func(t *testing.T) {
		b.LazyGetWithContext(context.TODO(), "post")

		assert.Equal(t, bb.HandleLazyGetCalls, 1)
	})
	t.Run("not providing context", func(t *testing.T) {
		b.LazyGet("post")

		assert.Equal(t, bb.HandleLazyGetCalls, 2)
	})
}

func TestBox_ListFromPeriodHasForwardedToHandleListFromPeriod(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	t.Run("providing context", func(t *testing.T) {
		b.ListFromPeriodWithContext(context.TODO(), time.Now(), time.Now(), 1)

		assert.Equal(t, bb.HandleListFromPeriodCalls, 1)
	})
	t.Run("not providing context", func(t *testing.T) {
		b.ListFromPeriod(time.Now(), time.Now(), 1)

		assert.Equal(t, bb.HandleListFromPeriodCalls, 2)
	})
}

func TestBox_DeleteHasForwardedToHandleDelete(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	t.Run("providing context", func(t *testing.T) {
		b.DeleteWithContext(context.TODO(), "post")

		assert.Equal(t, bb.HandleDeleteCalls, 1)
	})
	t.Run("not providing context", func(t *testing.T) {
		b.Delete("post")

		assert.Equal(t, bb.HandleDeleteCalls, 2)
	})
}

func TestBox_CleanHasForwardedToHandleClean(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	t.Run("providing context", func(t *testing.T) {
		b.CleanWithContext(context.TODO())

		assert.Equal(t, bb.HandleCleanCalls, 1)
	})
	t.Run("not providing context", func(t *testing.T) {
		b.Clean()

		assert.Equal(t, bb.HandleCleanCalls, 2)
	})
}
