package filesystem

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func createBoxFolder(t *testing.T, b *box) {
	if err := os.MkdirAll(filepath.Join(b.p.path, b.id), 0666); err != nil {
		t.Fatal("unable to create box folder")
	}
}

func assertContentFileIsCreated(t *testing.T, b *box, id string) {
	t.Helper()

	_, err := os.Stat(filepath.Join(b.p.path, b.id, id))
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		t.Fatal("didn't create box content file")
		return
	}
	log.Fatal("unable to check if box content file is created")
}

func TestBox_Post(t *testing.T) {
	id := "box_1"
	p := &provider{"Mailbox"}
	b := &box{p, id}
	createBoxFolder(t, b)

	t.Run("post content and return its identifier", func(t *testing.T) {
		err := b.Post("1", "foo")

		assert.Nil(t, err)
		assertContentFileIsCreated(t, b, "1")
	})

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(p.path)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}
