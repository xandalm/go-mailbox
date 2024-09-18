package memory

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestBox_Post(t *testing.T) {

	t.Run("post content", func(t *testing.T) {
		b := &box{
			contents: map[string]any{},
		}

		err := b.Post("1", "lorem ipsum")

		assert.Nil(t, err)

		assert.NotEmpty(t, b.contents)
	})
	t.Run("returns error because id duplication", func(t *testing.T) {
		b := &box{
			contents: map[string]any{"1": "foo"},
		}

		err := b.Post("1", "bar")
		assert.Error(t, err, ErrRepeatedContentIdentifier)
	})
	t.Run("returns error because nil content", func(t *testing.T) {
		b := &box{
			contents: map[string]any{},
		}

		err := b.Post("1", nil)
		assert.Error(t, err, ErrPostingNilContent)
	})
}

func TestBox_Get(t *testing.T) {
	t.Run("returns the content by post identifier", func(t *testing.T) {
		b := &box{
			contents: map[string]any{"1": "foo"},
		}

		got, err := b.Get("1")

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
			contents: map[string]any{"1": "foo"},
		}

		err := b.Delete("1")

		assert.Nil(t, err)
		assert.Empty(t, b.contents)
	})
}

func TestBox_Clean(t *testing.T) {
	t.Run("remove all content", func(t *testing.T) {
		b := &box{
			contents: map[string]any{"1": "foo", "2": "bar", "3": struct{ data any }{"foobarbaz"}},
		}

		err := b.Clean()

		assert.Nil(t, err)
		assert.Empty(t, b.contents)
	})
}
