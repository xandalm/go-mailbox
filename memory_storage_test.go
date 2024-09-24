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

	st.boxes["box_2"] = memoryStorageBox{}

	t.Run("returns error because id already exists", func(t *testing.T) {

		err := st.CreateBox("box_2")

		assert.Error(t, err, ErrRepeatedBoxIdentifier)
	})
}

func TestMemoryStorage_ListingBoxes(t *testing.T) {
	st := &MemoryStorage{
		boxes: map[string]memoryStorageBox{
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
		boxes: map[string]memoryStorageBox{
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
		boxes: map[string]memoryStorageBox{
			"box_1": {
				content: map[string]Bytes{
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
