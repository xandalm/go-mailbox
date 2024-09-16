package memory

import (
	"testing"

	"github.com/xandalm/go-mailbox"
	"github.com/xandalm/go-testing/assert"
)

func TestBox_Post(t *testing.T) {

	t.Run("post content and return its identifier", func(t *testing.T) {
		b := &box{
			contents: map[any]any{},
		}

		got, err := b.Post("lorem ipsum")

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.NotEmpty(t, b.contents)
	})
	t.Run("don't repeat identifier", func(t *testing.T) {
		b := &box{
			contents: map[any]any{},
		}

		id1, err := b.Post("foo")
		assert.Nil(t, err)
		assert.NotEmpty(t, b.contents)

		id2, err := b.Post("bar")
		assert.Nil(t, err)
		assert.NotEmpty(t, b.contents)

		assert.NotEqual(t, id1, id2)
	})
	t.Run("returns error because nil content", func(t *testing.T) {
		b := &box{
			contents: map[any]any{},
		}

		id, err := b.Post(nil)
		assert.Nil(t, id)
		assert.Error(t, err, mailbox.ErrNilContent)
	})
}

func TestBox_Get(t *testing.T) {
	t.Run("returns the content by post identifier", func(t *testing.T) {
		b := &box{
			contents: map[any]any{1: "foo"},
		}

		got, err := b.Get(1)

		assert.Nil(t, err)
		assert.NotNil(t, got)

		content, ok := got.(string)

		if !ok {
			t.Fatal("didn't get the expected content type")
		}

		assert.Equal(t, content, "foo")
	})
}

func TestBox_Delete(t *testing.T) {
	t.Run("remove content", func(t *testing.T) {
		b := &box{
			contents: map[any]any{1: "foo"},
		}

		err := b.Delete(1)

		assert.Nil(t, err)
		assert.Empty(t, b.contents)
	})
}
