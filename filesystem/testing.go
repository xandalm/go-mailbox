package filesystem

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func newCleanUpFunc(p *provider) func() {
	return func() {
		if p == nil {
			return
		}
		if p.f != nil {
			p.f.Close()
			for _, b := range p.boxes {
				if b.f != nil {
					b.f.Close()
				}
			}
		}
		if err := os.RemoveAll(filepath.Join(p.path)); err != nil {
			log.Fatalf("unable to remove residual data, %v", err)
		}
	}
}

func createFolder(path, dir string) {
	if err := os.MkdirAll(filepath.Join(path, dir), 0666); err != nil {
		log.Fatalf("unable to create folder, %v", err)
	}
}

func createProvider(path, dir string) *provider {
	createFolder(path, dir)
	path = filepath.Join(path, dir)
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("unable to open provider file, %v", err)
	}
	p := &provider{
		sync.RWMutex{},
		f,
		[]*box{},
		path,
	}
	return p
}

func createBox(p *provider, id string) *box {
	createFolder(p.path, id)
	path := filepath.Join(p.path, id)
	pos, _ := slices.BinarySearchFunc(p.boxes, id, func(b *box, id string) int {
		return strings.Compare(b.id, id)
	})
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("unable to open box file, %v", err)
	}
	b := &box{
		sync.RWMutex{},
		f,
		p,
		id,
	}
	p.boxes = slices.Insert(p.boxes, pos, b)
	return b
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
