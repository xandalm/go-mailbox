package mailbox

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestBoxRequesting(t *testing.T) {

	p := &stubProvider{}
	m := &manager{
		p: p,
	}

	got, err := m.RequestBox("box_1")

	assert.Nil(t, err)
	assert.NotNil(t, got)

	t.Run("returns error for failing provider", func(t *testing.T) {
		p := &stubFailingProvider{}
		m := NewManager(p)

		_, got := m.RequestBox("box_1")

		assert.NotNil(t, got)
	})
}

func TestBoxErasing(t *testing.T) {

	p := &stubProvider{
		Boxes: []*stubBox{{"box_1"}},
	}
	m := &manager{
		p: p,
	}

	err := m.EraseBox("box_1")

	assert.Nil(t, err)

	t.Run("returns error from provider", func(t *testing.T) {
		m := &manager{&stubFailingProvider{}}
		err := m.EraseBox("box_1")

		assert.NotNil(t, err)
	})
}

func TestCheckingForBox(t *testing.T) {

	provider := &stubProvider{
		Boxes: []*stubBox{{"box_1"}},
	}
	manager := NewManager(provider)

	t.Run("returns true", func(t *testing.T) {
		got, err := manager.ContainsBox("box_1")
		assert.Nil(t, err)
		assert.True(t, got)
	})
	t.Run("returns false", func(t *testing.T) {
		got, err := manager.ContainsBox("box_2")
		assert.Nil(t, err)
		assert.False(t, got)
	})
}
