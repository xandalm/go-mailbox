package filesystem

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func createBoxFolder(t *testing.T, b *box) {
	if err := os.MkdirAll(filepath.Join(b.p.path, b.id), 0666); err != nil {
		t.Errorf("unable to create box folder, %v", err)
	}
}

func createBoxContentFile(t *testing.T, b *box, id string, content Bytes) {
	filename := filepath.Join(b.p.path, b.id, id)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		t.Errorf("unable to create box content file, %v", err)
	}
	defer f.Close()
	f.Write(content)
}

func isContentFileCreated(b *box, id string) bool {

	_, err := os.Stat(filepath.Join(b.p.path, b.id, id))
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Fatal("unable to check if box content file is created")
	return false
}

func assertContentFileIsCreated(t *testing.T, b *box, id string) {
	t.Helper()

	if !isContentFileCreated(b, id) {
		t.Fatalf("the box folder didn't have content file with id=%q", id)
	}
}

func assertContentFileHasData(t *testing.T, b *box, id string, content Bytes) {
	t.Helper()

	if !isContentFileCreated(b, id) {
		t.Fatalf("the box folder didn't have content file with id=%q", id)
	}
	f, err := os.Open(filepath.Join(b.p.path, b.id, id))
	if err != nil {
		log.Fatal("unable to open box content file")
	}
	defer f.Close()
	data, _ := io.ReadAll(f)
	assert.Equal(t, data, content)
}

func TestBox_Post(t *testing.T) {
	id := "box_1"
	p := &provider{"Mailbox"}
	b := &box{&rwImpl{}, p, id}
	createBoxFolder(t, b)

	t.Run("post content", func(t *testing.T) {
		content := Bytes("foo")
		err := b.Post("1", content)

		assert.Nil(t, err)
		assertContentFileHasData(t, b, "1", content)
	})

	t.Run("returns error because id duplication", func(t *testing.T) {
		assertContentFileIsCreated(t, b, "1")

		err := b.Post("1", Bytes("bar"))
		assert.Error(t, err, ErrRepeatedContentIdentifier)
	})

	t.Run("returns error because nil content", func(t *testing.T) {
		err := b.Post("2", nil)
		assert.Error(t, err, ErrPostingNilContent)
	})

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(p.path)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}

func TestBox_Get(t *testing.T) {
	id := "box_1"
	p := &provider{"Mailbox"}
	b := &box{&rwImpl{}, p, id}
	createBoxFolder(t, b)
	createBoxContentFile(t, b, "1", Bytes("foo"))

	t.Run("returns the content by post identifier", func(t *testing.T) {
		got, err := b.Get("1")

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got, Bytes("foo"))
	})

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(p.path)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}
