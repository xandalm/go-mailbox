package filesystem

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/xandalm/go-mailbox"
	"github.com/xandalm/go-testing/assert"
)

var errFoo = errors.New("some error")

type stubFailingFsHandler struct{}

func (rw *stubFailingFsHandler) Read(path, id string) ([]byte, error) {
	return nil, errFoo
}

func (rw *stubFailingFsHandler) Write(path, id string, data []byte) error {
	return errFoo
}

func (rw *stubFailingFsHandler) Delete(path, id string) error {
	return errFoo
}

func (rw *stubFailingFsHandler) Clean(path string) error {
	return errFoo
}

func createBoxContentFile(b *box, id string, content Bytes) {
	filename := filepath.Join(b.p.path, b.id, id)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("unable to create box content file, %v", err)
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

func assertContentFileExists(t *testing.T, b *box, id string) {
	t.Helper()

	if !isContentFileCreated(b, id) {
		t.Fatalf("the box folder didn't have content file with id=%q", id)
	}
}

func assertContentFileNotExists(t *testing.T, b *box, id string) {
	t.Helper()

	if isContentFileCreated(b, id) {
		t.Fatalf("the box folder has the content file with id=%q", id)
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
	path := t.TempDir()
	dir := "Mailbox"
	p := createProvider(path, dir)
	id := "box_1"
	b := createBox(p, id)

	t.Run("post content", func(t *testing.T) {
		content := Bytes("foo")
		err := b.Post("1", content)

		assert.Nil(t, err)
		assertContentFileHasData(t, b, "1", content)
	})

	t.Run("returns error because id duplication", func(t *testing.T) {
		assertContentFileExists(t, b, "1")

		err := b.Post("1", Bytes("bar"))
		assert.Error(t, err, ErrRepeatedContentIdentifier)
	})

	t.Run("returns error because nil content", func(t *testing.T) {
		err := b.Post("2", nil)
		assert.Error(t, err, ErrPostingNilContent)
	})

	t.Run("returns error because unexpected condition", func(t *testing.T) {
		b.fs = &stubFailingFsHandler{}

		err := b.Post("2", Bytes("bar"))
		assert.Error(t, err, mailbox.ErrUnableToPostContent)
	})

	t.Cleanup(newCleanUpFunc(p))
}

func TestBox_Get(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"
	p := createProvider(path, dir)
	id := "box_1"
	b := createBox(p, id)
	createBoxContentFile(b, "1", Bytes("foo"))

	t.Run("returns the content by post identifier", func(t *testing.T) {
		got, err := b.Get("1")

		assert.Nil(t, err)
		assert.NotNil(t, got)

		assert.Equal(t, got, Bytes("foo"))
	})

	t.Run("returns error because post file don't exist", func(t *testing.T) {
		data, err := b.Get("2")

		assert.Nil(t, data)
		assert.Error(t, err, ErrContentNotFound)
	})

	t.Run("returns error because unexpected condition", func(t *testing.T) {
		b.fs = &stubFailingFsHandler{}

		_, err := b.Get("1")
		assert.Error(t, err, mailbox.ErrUnableToReadContent)
	})

	t.Cleanup(newCleanUpFunc(p))
}

func TestBox_Delete(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"
	p := createProvider(path, dir)
	id := "box_1"
	b := createBox(p, id)
	createBoxContentFile(b, "1", Bytes("foo"))

	t.Run("delete content", func(t *testing.T) {
		err := b.Delete("1")

		assert.Nil(t, err)
		assertContentFileNotExists(t, b, "1")
	})

	t.Run("returns error because unexpected condition", func(t *testing.T) {
		createBoxContentFile(b, "1", Bytes("foo"))
		b.fs = &stubFailingFsHandler{}
		err := b.Delete("1")

		assert.Error(t, err, mailbox.ErrUnableToDeleteContent)
		assertContentFileExists(t, b, "1")
	})

	t.Run("do nothing on not found content", func(t *testing.T) {
		b.fs = &fsHandlerImpl{}
		err := b.Delete("2")

		assert.Nil(t, err)
		assertContentFileNotExists(t, b, "2")
	})

	t.Cleanup(newCleanUpFunc(p))
}

func TestBox_Clean(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"
	p := createProvider(path, dir)
	id := "box_1"
	b := createBox(p, id)
	createBoxContentFile(b, "1", Bytes("foo"))
	createBoxContentFile(b, "2", Bytes("bar"))

	t.Run("remove all content", func(t *testing.T) {
		err := b.Clean()

		assert.Nil(t, err)
		assertContentFileNotExists(t, b, "1")
		assertContentFileNotExists(t, b, "2")
	})

	t.Run("returns error because unexpected condition", func(t *testing.T) {
		b.fs = &stubFailingFsHandler{}
		err := b.Clean()

		assert.Error(t, err, mailbox.ErrUnableToCleanBox)
	})

	t.Cleanup(newCleanUpFunc(p))
}
