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

func TestMemoryStorage_Listing(t *testing.T) {
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
