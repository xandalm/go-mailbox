package mailbox

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestBoxRequesting(t *testing.T) {

	p := &stubProvider{}
	m := &manager{
		p:   p,
		idx: []string{},
	}

	got, err := m.RequestBox("box_1")

	assert.Nil(t, err)
	assert.NotNil(t, got)

	assert.Contains(t, m.idx, "box_1")
	assert.ContainsFunc(t, p.Boxes, "box_1", func(e *stubBox, lf string) bool {
		return e.Id == lf
	})

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
		p:   p,
		idx: []string{"box_1"},
	}

	err := m.EraseBox("box_1")

	assert.Nil(t, err)
	assert.NotContains(t, m.idx, "box_1")
	assert.NotContainsFunc(t, p.Boxes, "box_1", func(sb *stubBox, s string) bool {
		return sb.Id == s
	})

	t.Run("returns error for inexistent box", func(t *testing.T) {
		err := m.EraseBox("box_2")

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
