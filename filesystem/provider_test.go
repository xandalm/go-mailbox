package filesystem

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestNewProvider(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"

	got := NewProvider(path, dir).(*provider)

	assert.NotNil(t, got)
	if got.f == nil {
		t.Error("didn't open dir")
	}
	wantPath := filepath.Join(path, dir)
	if got.path != wantPath {
		t.Errorf("got provider path %s, but want %s", got.path, wantPath)
	}

	t.Run("load existing boxes", func(t *testing.T) {
		dir := "FilledMailbox"

		createFolder(path, dir)
		pDirPath := filepath.Join(path, dir)
		createFolder(pDirPath, "box_1")
		createFolder(pDirPath, "box_2")

		got := NewProvider(path, dir).(*provider)

		assert.NotNil(t, got)
		if got.f == nil {
			t.Error("didn't open dir")
		}
		wantPath := filepath.Join(path, dir)
		if got.path != wantPath {
			t.Errorf("got provider path %s, but want %s", got.path, wantPath)
		}

		assert.ContainsFunc(t, got.boxes, "box_1", func(b *boxFile, id string) bool {
			return b.id == id
		})
		assert.ContainsFunc(t, got.boxes, "box_2", func(b *boxFile, id string) bool {
			return b.id == id
		})

		t.Cleanup(newCleanUpFunc(got))
	})

	t.Cleanup(newCleanUpFunc(got))
}

func TestProvider_Create(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"
	p := createProvider(path, dir)

	t.Run("create and return box", func(t *testing.T) {
		got, err := p.Create("box_1")

		assert.Nil(t, err)
		assert.NotNil(t, got)
		b := got.(*box)
		assert.Equal(t, b.p, p)
		assert.Equal(t, b.bf.id, "box_1")
		assert.NotNil(t, b.bf.f)

		assert.ContainsFunc(t, p.boxes, "box_1", func(b *boxFile, id string) bool {
			return b.id == id
		})
		entry, osErr := os.ReadDir(filepath.Join(p.path))
		if osErr != nil {
			t.Fatalf("unable to check dir, %v", osErr)
		}
		assert.ContainsFunc(t, entry, "box_1", func(de fs.DirEntry, s string) bool {
			return de.Name() == s
		})
	})

	t.Run("return error by empty id", func(t *testing.T) {
		b, got := p.Create("")

		assert.Nil(t, b)
		assert.Error(t, got, ErrEmptyBoxIdentifier)
	})

	t.Run("returns error for the id duplicity", func(t *testing.T) {
		b, got := p.Create("box_1")

		assert.Nil(t, b)
		assert.Error(t, got, ErrRepeatedBoxIdentifier)
	})

	t.Cleanup(newCleanUpFunc(p))
}

func TestProvider_Get(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"

	p := createProvider(path, dir)
	createBox(p, "box_1")

	t.Run("return box", func(t *testing.T) {
		got, err := p.Get("box_1")

		assert.Nil(t, err)
		assert.NotNil(t, got)
		b := got.(*box)
		assert.Equal(t, b.p, p)
		assert.Equal(t, b.bf.id, "box_1")
		assert.NotNil(t, b.bf.f)
	})

	t.Run("return error because box doesn't exist", func(t *testing.T) {
		b, got := p.Get("box_2")

		assert.Nil(t, b)
		assert.Error(t, got, ErrBoxNotFound)
	})

	t.Cleanup(newCleanUpFunc(p))
}

func TestProvider_Contains(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"

	p := createProvider(path, dir)
	createBox(p, "box_1")

	t.Run("returns true and nil error", func(t *testing.T) {
		got, err := p.Contains("box_1")

		assert.Nil(t, err)
		assert.True(t, got)
	})

	t.Run("returns false and nil error", func(t *testing.T) {
		got, err := p.Contains("box_2")

		assert.Nil(t, err)
		assert.False(t, got)
	})

	t.Cleanup(newCleanUpFunc(p))
}

func TestProvider_Delete(t *testing.T) {
	path := t.TempDir()
	dir := "Mailbox"

	p := createProvider(path, dir)
	createBox(p, "box_1")

	t.Run("delete box", func(t *testing.T) {
		got := p.Delete("box_1")

		assert.Nil(t, got)
		assert.NotContainsFunc(t, p.boxes, "box_1", func(b *boxFile, id string) bool {
			return b.id == id
		})
		if _, err := os.Stat(filepath.Join(path, dir, "box_1")); err == nil {
			t.Error("didn't delete box folder")
		} else if !os.IsNotExist(err) {
			log.Fatal("unable to check box folder existence")
		}
	})

	t.Cleanup(newCleanUpFunc(p))
}
