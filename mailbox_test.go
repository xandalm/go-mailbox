package mailbox

import (
	"testing"

	"github.com/xandalm/go-session/testing/assert"
)

func TestBoxCreating(t *testing.T) {

	storage := &stubStorage{[]Box{}}
	manager := NewManager(storage)

	got, err := manager.CreateBox("box_1")
	want := Box{
		Id: "box_1",
	}

	assert.NoError(t, err)

	assert.Equal(t, got, want)

	assert.NotEmpty(t, storage.collec)

	assert.Equal(t, storage.collec[0], want)

	t.Run("returns error for duplicity", func(t *testing.T) {
		box, got := manager.CreateBox("box_1")
		want := ErrBoxIDDuplicity

		assert.Equal(t, box, Box{})
		assert.Error(t, got, want)
	})
}

type stubStorage struct {
	collec []Box
}

func (s *stubStorage) Save(b Box) error {
	s.collec = append(s.collec, b)
	return nil
}

func (s *stubStorage) List() ([]string, error) {
	return nil, nil
}
