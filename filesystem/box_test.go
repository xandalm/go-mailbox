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

func createBoxContentFile(t *testing.T, b *box, id, content string) {
	filename := filepath.Join(b.p.path, b.id, id)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		t.Errorf("unable to create box content file, %v", err)
	}
	defer f.Close()
	f.Write([]byte(content))
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

func assertContentFileHasData(t *testing.T, b *box, id string, content string) {
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
	assert.Equal(t, string(data), content)
}

func TestBox_Post(t *testing.T) {
	id := "box_1"
	p := &provider{"Mailbox"}
	b := &box{p, id}
	createBoxFolder(t, b)

	t.Run("post content", func(t *testing.T) {
		err := b.Post("1", "foo")

		assert.Nil(t, err)
		assertContentFileHasData(t, b, "1", "foo")
	})

	t.Run("returns error because id duplication", func(t *testing.T) {
		assertContentFileIsCreated(t, b, "1")

		err := b.Post("1", "bar")
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
	b := &box{p, id}
	createBoxFolder(t, b)
	createBoxContentFile(t, b, "1", "foo")

	t.Run("returns the content by post identifier", func(t *testing.T) {
		got, err := b.Get("1")

		assert.Nil(t, err)
		assert.NotNil(t, got)

		content, ok := got.(string)

		if !ok {
			t.Fatal("didn't get the expected content type")
		}

		assert.Equal(t, content, "foo")
	})

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(p.path)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}
