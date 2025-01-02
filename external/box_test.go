package external

import (
	"context"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestNewBox(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	assert.NotNil(t, b)
}

var dummyContent = []byte("abc")

func TestBox_PostHasForwardedToOnPost(t *testing.T) {
	bb := &spyBoxBridge{}
	b := NewBox(bb)

	t.Run("providing context", func(t *testing.T) {
		b.PostWithContext(context.TODO(), dummyContent)

		assert.Equal(t, bb.OnPostCalls, 1)
	})
	t.Run("not providing context", func(t *testing.T) {
		b.Post(dummyContent)

		assert.Equal(t, bb.OnPostCalls, 2)
	})
}
