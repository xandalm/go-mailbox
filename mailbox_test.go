package mailbox

import (
	"testing"

	"github.com/xandalm/go-session/testing/assert"
)

func TestCreateBox(t *testing.T) {

	got, err := CreateBox("box_1")
	want := Box{
		Id: "box_1",
	}

	assert.NoError(t, err)

	assert.Equal(t, got, want)
}
