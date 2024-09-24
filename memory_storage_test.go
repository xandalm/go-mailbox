package mailbox

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestMemoryStorage_CreatingBox(t *testing.T) {
	st := NewMemoryStorage()

	t.Run("create box in storage", func(t *testing.T) {
		err := st.CreateBox("box_1")

		assert.Nil(t, err)
		assert.NotEmpty(t, st.boxes)

		if _, ok := st.boxes["box_1"]; !ok {
			t.Errorf("didn't create box in storage")
		}
	})

	st.boxes["box_2"] = &memoryStorageBox{}

	t.Run("returns error because id already exists", func(t *testing.T) {

		err := st.CreateBox("box_2")

		assert.Error(t, err, ErrRepeatedBoxIdentifier)
	})
}

func TestMemoryStorage_ListingBoxes(t *testing.T) {
	st := &MemoryStorage{
		boxes: map[string]*memoryStorageBox{
			"box_1": {},
			"box_2": {},
		},
	}

	t.Run("returns all boxes ids", func(t *testing.T) {
		ids, err := st.ListBoxes()

		assert.Nil(t, err)
		assert.NotNil(t, ids)
		assert.NotEmpty(t, ids)

		assert.Contains(t, ids, "box_1")
		assert.Contains(t, ids, "box_2")
	})
}

func TestMemoryStorage_DeletingBox(t *testing.T) {
	st := &MemoryStorage{
		boxes: map[string]*memoryStorageBox{
			"box_1": {},
			"box_2": {},
		},
	}

	id := "box_1"
	err := st.DeleteBox(id)

	assert.Nil(t, err)

	if _, ok := st.boxes[id]; ok {
		t.Errorf("storage still contains %s", id)
	}
}

func TestMemoryStorage_CleanBox(t *testing.T) {
	st := &MemoryStorage{
		boxes: map[string]*memoryStorageBox{
			"box_1": {
				content: map[string][]byte{
					"a7da5": Bytes("foo"),
				},
			},
		},
	}

	err := st.CleanBox("box_1")

	assert.Nil(t, err)
	box, ok := st.boxes["box_1"]
	assert.True(t, ok, "the box must be kept in %v", st.boxes)
	assert.Empty(t, box.content)
}

func TestMemoryStorage_CreateContent(t *testing.T) {
	st := &MemoryStorage{
		boxes: map[string]*memoryStorageBox{
			"box_1": {
				content: map[string][]byte{
					"a7da5": []byte("foo"),
				},
			},
		},
	}

	t.Run("create box content in storage", func(t *testing.T) {
		err := st.CreateContent("box_1", "b042e", Bytes("bar"))

		assert.Nil(t, err)
		box := st.boxes["box_1"]
		assert.NotNil(t, box.content)
		content, ok := box.content["b042e"]
		assert.True(t, ok, "the content isn't created in %v", box.content)
		assert.Equal(t, content, Bytes("bar"))
	})

	t.Run("returns error because the content id already exists", func(t *testing.T) {
		err := st.CreateContent("box_1", "a7da5", Bytes("baz"))

		assert.Error(t, err, ErrRepeatedContentIdentifier)
	})

	t.Run("returns error because box doesn't exist", func(t *testing.T) {
		err := st.CreateContent("box_2", "a7da5", Bytes("baz"))

		assert.Error(t, err, ErrBoxNotFoundToPostContent)
	})
}

func TestMemoryStorage_ReadContent(t *testing.T) {
	st := &MemoryStorage{
		boxes: map[string]*memoryStorageBox{
			"box_1": {
				content: map[string][]byte{
					"a7da5": []byte("foo"),
				},
			},
		},
	}

	t.Run("read box content from storage", func(t *testing.T) {
		got, err := st.ReadContent("box_1", "a7da5")
		want := []byte("foo")

		assert.Nil(t, err)
		assert.Equal(t, got, want)
	})

	t.Run("returns error because box doesn't exist", func(t *testing.T) {
		content, err := st.ReadContent("box_2", "a7da5")

		assert.Nil(t, content)
		assert.Error(t, err, ErrBoxNotFoundToReadContent)
	})
}

func TestMemoryStorage_DeleteContent(t *testing.T) {
	st := &MemoryStorage{
		boxes: map[string]*memoryStorageBox{
			"box_1": {
				content: map[string][]byte{
					"a7da5": []byte("foo"),
				},
			},
		},
	}

	err := st.DeleteContent("box_1", "a7da5")

	assert.Nil(t, err)
	box := st.boxes["box_1"]
	assert.NotNil(t, box)
	assert.NotNil(t, box.content)
	if _, ok := box.content["a7da5"]; ok {
		t.Errorf("didn't delete the box content")
	}
}
