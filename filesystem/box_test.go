package filesystem

import (
	"bytes"
	"testing"
	"time"

	"github.com/xandalm/go-mailbox"
	"github.com/xandalm/go-testing/assert"
)

func TestBox_Post(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"
	p := createProvider(path, dir)
	id := "box_1"
	b := createBox(p, id)

	t.Run("post content", func(t *testing.T) {
		content := Bytes("foo")
		ct, err := b.Post("1", content)

		assert.Nil(t, err)
		assert.NotZero(t, ct)
		assertContentFileHasData(t, b, "1", content)
	})

	t.Run("returns error because id duplication", func(t *testing.T) {
		assertContentFileExists(t, b, "1")

		ct, err := b.Post("1", Bytes("bar"))
		assert.Zero(t, ct)
		assert.Error(t, err, ErrRepeatedContentIdentifier)
	})

	t.Run("returns error because nil content", func(t *testing.T) {
		ct, err := b.Post("2", nil)
		assert.Zero(t, ct)
		assert.Error(t, err, ErrPostingNilContent)
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

		assert.Equal(t, got, mailbox.Data{Content: Bytes("foo")})
	})

	t.Run("returns error because post file don't exist", func(t *testing.T) {
		data, err := b.Get("2")

		assert.Zero(t, data)
		assert.Error(t, err, ErrContentNotFound)
	})

	t.Cleanup(newCleanUpFunc(p))
}

func TestBox_GetFromPeriod(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"
	p := createProvider(path, dir)
	id := "box_1"
	b := createBox(p, id)
	createBoxContentFile(b, "b2f07", Bytes("foo"))
	time.Sleep(time.Millisecond)
	createBoxContentFile(b, "f348c", Bytes("bar"))
	time.Sleep(time.Millisecond)
	createBoxContentFile(b, "a5c01", Bytes("baz"))

	end := time.Now().UnixNano()
	begin := end - int64(2*time.Millisecond)

	got, err := b.GetFromPeriod(begin, end)

	assert.Nil(t, err)
	assert.NotNil(t, got)

	assert.NotEmpty(t, got)
	if len(got) != 2 {
		t.Fatalf("didn't get expected length, got %d want 2", len(got))
	}

	assert.EqualFunc(t, got[0], mailbox.Data{Content: Bytes("bar")}, func(a, b mailbox.Data) bool {
		return bytes.Equal(a.Content, b.Content)
	})
	assert.EqualFunc(t, got[1], mailbox.Data{Content: Bytes("baz")}, func(a, b mailbox.Data) bool {
		return bytes.Equal(a.Content, b.Content)
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

	t.Run("do nothing on not found content", func(t *testing.T) {
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

	t.Cleanup(newCleanUpFunc(p))
}
