package filesystem

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestProvider_Create(t *testing.T) {
	path := ""
	dir := "Mailbox"
	p := NewProvider(path, dir)

	t.Run("create and return box", func(t *testing.T) {
		got, err := p.Create("box_1")

		assert.Nil(t, err, "expected nil but got %v", err)
		assert.NotNil(t, got)
		assert.Equal(t, *got.(*box), box{id: "box_1"})

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

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(path, dir)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}

func TestProvider_Get(t *testing.T) {
	path := ""
	dir := "Mailbox"
	p := NewProvider(path, dir)

	p.Create("box_1")

	t.Run("return box", func(t *testing.T) {
		got, err := p.Get("box_1")

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, *got.(*box), box{id: "box_1"})
	})

	t.Run("return error because box doesn't exist", func(t *testing.T) {
		b, got := p.Get("box_2")

		assert.Nil(t, b)
		assert.Error(t, got, ErrBoxNotFound)
	})

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(path, dir)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}

func TestProvider_Delete(t *testing.T) {
	path := ""
	dir := "Mailbox"
	p := NewProvider(path, dir)

	p.Create("box_1")

	t.Run("delete box", func(t *testing.T) {
		got := p.Delete("box_1")

		assert.Nil(t, got)
		if _, err := os.Stat(filepath.Join(path, dir, "box_1")); err == nil {
			t.Error("didn't delete box folder")
		} else if !os.IsNotExist(err) {
			log.Fatal("unable to check box folder existence")
		}
	})

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(path, dir)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}

func TestProvider_List(t *testing.T) {
	path := ""
	dir := "Mailbox"
	p := NewProvider(path, dir)

	p.Create("box_1")
	p.Create("box_2")

	t.Run("return boxes id list", func(t *testing.T) {
		got, err := p.List()

		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Contains(t, got, "box_1")
		assert.Contains(t, got, "box_2")
	})

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(path, dir)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}
