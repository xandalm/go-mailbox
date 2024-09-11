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

	t.Run("returns error for provider failing while creating box", func(t *testing.T) {
		provider := &stubFailingProvider{}
		manager := NewManager(provider)

		_, got := manager.RequestBox("box_1")

		assert.AnError(t, got)
	})
}
