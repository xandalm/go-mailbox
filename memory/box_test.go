package memory

import (
	"container/list"
	"testing"
	"time"

	"github.com/xandalm/go-mailbox"
	"github.com/xandalm/go-testing/assert"
)

func TestBox_Post(t *testing.T) {

	t.Run("post content", func(t *testing.T) {
		id := "1"
		content := Bytes("lorem ipsum")
		b := &box{
			data:     list.New(),
			dataById: make(map[string]*list.Element),
		}

		ct, err := b.Post(id, content)

		assert.Nil(t, err)
		assert.NotZero(t, ct)

		if b.data.Len() == 0 {
			t.Fatal("expected not to be empty")
		}
		assert.NotEmpty(t, b.dataById)

		onMap, ok1 := b.dataById[id]
		if !ok1 {
			t.Fatalf("didn't have key %s in %v", id, b.dataById)
		}

		assert.Equal(t, onMap, b.data.Front())

		reg := onMap.Value.(*registry)
		wantReg := registry{id, ct, content}

		assert.Equal(t, *reg, wantReg)
	})
	t.Run("returns error because id duplication", func(t *testing.T) {
		b := &box{
			data:     list.New(),
			dataById: map[string]*list.Element{},
		}

		reg := &registry{"1", time.Now().UnixNano(), Bytes("foo")}
		b.dataById[reg.id] = b.data.PushBack(reg)

		ct, err := b.Post("1", Bytes("bar"))
		assert.Zero(t, ct)
		assert.Error(t, err, ErrRepeatedContentIdentifier)
	})
	t.Run("returns error because nil content", func(t *testing.T) {
		b := &box{
			data:     list.New(),
			dataById: map[string]*list.Element{},
		}

		ct, err := b.Post("1", nil)
		assert.Zero(t, ct)
		assert.Error(t, err, ErrPostingNilContent)
	})
}

func TestBox_Get(t *testing.T) {
	t.Run("returns the content by post identifier", func(t *testing.T) {
		b := &box{
			data:     list.New(),
			dataById: map[string]*list.Element{},
		}

		reg := &registry{"1", time.Now().UnixNano(), Bytes("foo")}
		b.dataById[reg.id] = b.data.PushBack(reg)

		got, err := b.Get("1")
		want := mailbox.Data{
			CreationTime: reg.ct,
			Content:      reg.c,
		}

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got, want)
	})
}

func TestBox_GetFromPeriod(t *testing.T) {
	t.Run("returns the content by post identifier", func(t *testing.T) {
		b := &box{
			data:     list.New(),
			dataById: map[string]*list.Element{},
		}

		now := time.Now().UnixNano()

		ct1 := now - int64(2*time.Second)
		reg1 := &registry{"1", ct1, Bytes("foo")}
		b.dataById[reg1.id] = b.data.PushBack(reg1)

		ct2 := now - int64(time.Second)
		reg2 := &registry{"2", ct2, Bytes("bar")}
		b.dataById[reg2.id] = b.data.PushBack(reg2)

		ct3 := now
		reg3 := &registry{"3", ct3, Bytes("baz")}
		b.dataById[reg3.id] = b.data.PushBack(reg3)

		got, err := b.GetFromPeriod(ct2, ct3)
		want := []mailbox.Data{
			{CreationTime: reg2.ct, Content: reg2.c},
			{CreationTime: reg3.ct, Content: reg3.c},
		}

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got, want)
	})
}

func TestBox_Delete(t *testing.T) {
	t.Run("remove content", func(t *testing.T) {
		b := &box{
			data:     list.New(),
			dataById: map[string]*list.Element{},
		}

		reg := &registry{"1", time.Now().UnixNano(), Bytes("foo")}
		b.dataById[reg.id] = b.data.PushBack(reg)

		err := b.Delete("1")

		assert.Nil(t, err)
		if b.data.Len() > 0 {
			t.Fatalf("expected empty but got %v", b.data)
		}
		assert.Empty(t, b.dataById)
	})
}

func TestBox_Clean(t *testing.T) {
	t.Run("remove all content", func(t *testing.T) {
		b := &box{
			data:     list.New(),
			dataById: map[string]*list.Element{},
		}

		reg1 := &registry{"1", time.Now().UnixNano(), Bytes("foo")}
		b.dataById[reg1.id] = b.data.PushBack(reg1)

		reg2 := &registry{"2", time.Now().UnixNano(), Bytes("bar")}
		b.dataById[reg2.id] = b.data.PushBack(reg2)

		reg3 := &registry{"3", time.Now().UnixNano(), Bytes("baz")}
		b.dataById[reg3.id] = b.data.PushBack(reg3)

		err := b.Clean()

		assert.Nil(t, err)
		if b.data.Len() > 0 {
			t.Fatalf("expected empty but got %v", b.data)
		}
		assert.Empty(t, b.dataById)
	})
}
