package mailbox

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestBoxRequesting(t *testing.T) {

	provider := &stubProvider{}
	manager := NewManager(provider)

	got, err := manager.RequestBox("box_1")

	assert.Nil(t, err)
	assert.NotNil(t, got)

	t.Run("returns error for failing provider", func(t *testing.T) {
		provider := &stubFailingProvider{}
		manager := NewManager(provider)

		_, got := manager.RequestBox("box_1")

		assert.NotNil(t, got)
	})
}

func TestBoxErasing(t *testing.T) {

	provider := &stubProvider{
		Boxes: []*stubBox{{"box_1"}},
	}
	manager := NewManager(provider)

	err := manager.EraseBox("box_1")

	assert.Nil(t, err)
	assert.Empty(t, provider.Boxes)

	t.Run("returns error for inexistent box", func(t *testing.T) {
		err := manager.EraseBox("box_2")

		assert.Error(t, err, ErrUnknownBox)
	})
}

func TestCheckingForBox(t *testing.T) {

	provider := &stubProvider{
		Boxes: []*stubBox{{"box_1"}},
	}
	manager := NewManager(provider)

	t.Run("returns true", func(t *testing.T) {
		got := manager.ContainsBox("box_1")
		assert.True(t, got)
	})
	t.Run("returns false", func(t *testing.T) {
		got := manager.ContainsBox("box_2")
		assert.False(t, got)
	})
}
