package memory

import (
	"testing"

	"github.com/xandalm/go-mailbox"
	"github.com/xandalm/go-session/testing/assert"
)

func TestProvider_Create(t *testing.T) {
	var p mailbox.Provider = &provider{
		boxes: map[string]*box{},
	}

	got, err := p.Create("box_1")

	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.NotEmpty(t, p.(*provider).boxes)

	t.Run("returns error for empty", func(t *testing.T) {
		b, got := p.Create("")

		assert.Nil(t, b)
		assert.Error(t, got, ErrEmptyBoxIdentifier)
	})
	t.Run("returns error for the id duplicity", func(t *testing.T) {
		b, got := p.Create("box_1")

		assert.Nil(t, b)
		assert.Error(t, got, ErrRepeatedBoxIdentifier)
	})
}

func TestProvider_Get(t *testing.T) {
	p := &provider{
		boxes: map[string]*box{
			"box_1": {},
		},
	}

	got, err := p.Get("box_1")
	want := p.boxes["box_1"]

	assert.NoError(t, err)
	assert.Equal(t, got.(*box), want)
}

func TestProvider_Delete(t *testing.T) {
	p := &provider{
		boxes: map[string]*box{
			"box_1": {},
		},
	}

	err := p.Delete("box_1")

	assert.NoError(t, err)
	assert.Empty(t, p.boxes)
}
