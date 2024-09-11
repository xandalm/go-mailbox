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

	t.Run("returns error for the id duplicity", func(t *testing.T) {
		b, got := p.Create("box_1")

		assert.Nil(t, b)
		assert.Error(t, got, mailbox.ErrBoxIDDuplicity)
	})
}
