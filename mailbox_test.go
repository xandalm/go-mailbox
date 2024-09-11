package mailbox

import (
	"testing"

	"github.com/xandalm/go-session/testing/assert"
)

func TestBoxCreating(t *testing.T) {

	provider := &stubProvider{}
	manager := NewManager(provider)

	got, err := manager.CreateBox("box_1")

	assert.NoError(t, err)
	assert.NotNil(t, got)

	t.Run("returns error for duplicity", func(t *testing.T) {
		box, got := manager.CreateBox("box_1")
		want := ErrBoxIDDuplicity

		assert.Nil(t, box)
		assert.Error(t, got, want)
	})
}

type stubProvider struct{}

func (s *stubProvider) Create(id string) (Box, error) {
	return nil, nil
}

func (s *stubProvider) List() ([]string, error) {
	return nil, nil
}
