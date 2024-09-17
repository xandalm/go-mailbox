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

	got, err := p.Create("box_1")

	assert.Nil(t, err)
	assert.NotNil(t, got)

	entry, osErr := os.ReadDir(filepath.Join(p.path))
	if osErr != nil {
		t.Fatalf("unable to check dir, %v", osErr)
	}
	assert.ContainsFunc(t, entry, "box_1", func(de fs.DirEntry, s string) bool {
		return de.Name() == s
	})

	t.Run("return error by empty id", func(t *testing.T) {
		b, got := p.Create("")

		assert.Nil(t, b)
		assert.Error(t, got, ErrEmptyBoxIdentifier)
	})

	t.Cleanup(func() {
		if err := os.RemoveAll(filepath.Join(path, dir)); err != nil {
			log.Fatal("unable to remove residual data")
		}
	})
}
