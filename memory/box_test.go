package memory

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestBox_Post(t *testing.T) {

	t.Run("post content", func(t *testing.T) {
		b := &box{
			contents: map[string]Bytes{},
		}

		err := b.Post("1", Bytes("lorem ipsum"))

		assert.Nil(t, err)

		assert.NotEmpty(t, b.contents)
	})
	t.Run("returns error because id duplication", func(t *testing.T) {
		b := &box{
			contents: map[string]Bytes{"1": Bytes("foo")},
		}

		err := b.Post("1", Bytes("bar"))
		assert.Error(t, err, ErrRepeatedContentIdentifier)
	})
	t.Run("returns error because nil content", func(t *testing.T) {
		b := &box{
			contents: map[string]Bytes{},
		}

		err := b.Post("1", nil)
		assert.Error(t, err, ErrPostingNilContent)
	})
}

func TestBox_Get(t *testing.T) {
	t.Run("returns the content by post identifier", func(t *testing.T) {
		content := Bytes("foo")
		b := &box{
			contents: map[string]Bytes{"1": content},
		}

		got, err := b.Get("1")

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got, content)
	})
}

func TestBox_Delete(t *testing.T) {
	t.Run("remove content", func(t *testing.T) {
		content := Bytes("foo")
		b := &box{
			contents: map[string]Bytes{"1": content},
		}

		err := b.Delete("1")

		assert.Nil(t, err)
		assert.Empty(t, b.contents)
	})
}

func TestBox_Clean(t *testing.T) {
	t.Run("remove all content", func(t *testing.T) {
		b := &box{
			contents: map[string]Bytes{
				"1": Bytes("foo"),
				"2": Bytes("bar"),
				"3": Bytes("baz"),
			},
		}

		err := b.Clean()

		assert.Nil(t, err)
		assert.Empty(t, b.contents)
	})
}
