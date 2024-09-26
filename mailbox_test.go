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

func TestPostingInBox(t *testing.T) {
	st := &stubStorage{
		Boxes: []stubStorageBox{{"box_1", make(map[string]Bytes)}},
	}
	b := &box{id: "box_1", st: st}

	t.Run("post content", func(t *testing.T) {
		cid := "data_1"
		data := []byte("foo")
		err := b.Post(cid, data)

		assert.Nil(t, err)
		if _, ok := st.Boxes[0].content[cid]; !ok {
			t.Fatalf("there's no %s in %v", cid, st.Boxes[0].content)
		}
	})

	t.Run("returns error because empty id", func(t *testing.T) {
		err := b.Post("", []byte("foo"))

		assert.Error(t, err, ErrEmptyContentIdentifier)
	})

	t.Run("returns erro because empty/nil data", func(t *testing.T) {

		assert.Error(t, b.Post("data_1", nil), ErrNothingToPost)

		assert.Error(t, b.Post("data_1", []byte{}), ErrNothingToPost)
	})
}

func TestReadingFromBox(t *testing.T) {
	st := &stubStorage{
		Boxes: []stubStorageBox{{"box_1", map[string]Bytes{"data_1": []byte("foo")}}},
	}
	b := &box{id: "box_1", st: st}

	t.Run("returns content data", func(t *testing.T) {
		data, err := b.Get("data_1")

		assert.Nil(t, err)
		assert.Equal(t, data, []byte("foo"))
	})
}
