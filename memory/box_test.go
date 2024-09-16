package memory

import (
	"testing"

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
		assert.NotNil(t, id1)
		assert.NotEmpty(t, b.contents)

		id2, err := b.Post("bar")
		assert.Nil(t, err)
		assert.NotNil(t, id2)
		assert.NotEmpty(t, b.contents)

		assert.NotEqual(t, id1, id2)
	})
}
