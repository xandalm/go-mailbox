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
}
