package mailbox

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestBoxRequesting(t *testing.T) {

	st := &stubStorage{}
	p := &provider{
		st:  st,
		idx: []string{},
	}

	got, err := p.RequestBox("box_1")

	assert.Nil(t, err)
	assert.NotNil(t, got)

	assert.Contains(t, p.idx, "box_1")
	assert.ContainsFunc(t, st.Boxes, "box_1", func(e stubStorageBox, lf string) bool {
		return e.id == lf
	})

	t.Run("returns error for failing provider", func(t *testing.T) {
		st := &stubFailingStorage{}
		m := NewProvider(st)

		_, got := m.RequestBox("box_1")

		assert.NotNil(t, got)
	})
}

func TestBoxErasing(t *testing.T) {

	st := &stubStorage{
		Boxes: []stubStorageBox{},
	}
	p := &provider{
		st:  st,
		idx: []string{"box_1"},
	}

	err := p.EraseBox("box_1")

	assert.Nil(t, err)
	assert.NotContains(t, p.idx, "box_1")
	assert.NotContainsFunc(t, st.Boxes, "box_1", func(sb stubStorageBox, s string) bool {
		return sb.id == s
	})

	t.Run("returns error for inexistent box", func(t *testing.T) {
		err := p.EraseBox("box_2")

		assert.Error(t, err, ErrUnknownBox)
	})
}

func TestCheckingForBox(t *testing.T) {

	st := &stubStorage{
		Boxes: []stubStorageBox{{"box_1", make(map[string]Bytes)}},
	}
	p := NewProvider(st)

	t.Run("returns true", func(t *testing.T) {
		got := p.ContainsBox("box_1")
		assert.True(t, got)
	})
	t.Run("returns false", func(t *testing.T) {
		got := p.ContainsBox("box_2")
		assert.False(t, got)
	})
}
