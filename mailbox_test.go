package mailbox

import (
	"testing"

	"github.com/xandalm/go-session/testing/assert"
)

func TestBoxRequesting(t *testing.T) {

	provider := &stubProvider{}
	manager := NewManager(provider)

	got, err := manager.RequestBox("box_1")

	assert.NoError(t, err)
	assert.NotNil(t, got)

	t.Run("returns error for failing provider", func(t *testing.T) {
		provider := &stubFailingProvider{}
		manager := NewManager(provider)

		_, got := manager.RequestBox("box_1")

		assert.AnError(t, got)
	})
}

func TestBoxErasing(t *testing.T) {

	provider := &stubProvider{
		Boxes: []*stubBox{{"box_1"}},
	}
	manager := NewManager(provider)

	err := manager.EraseBox("box_1")

	assert.NoError(t, err)
	assert.Empty(t, provider.Boxes)
}
