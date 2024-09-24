package mailbox

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func createStorageFolder(st *fileSystemStorage) {
	err := os.MkdirAll(st.path, 0666)
	if err != nil && !os.IsExist(err) {
		log.Fatal("unable to create storage folder")
	}
}

func createBoxFolder(st *fileSystemStorage, bid string) {
	if err := os.MkdirAll(filepath.Join(st.path, bid), 0666); err != nil {
		log.Fatalf("unable to create box folder, %v", err)
	}
}

func assertStorageFolderIsCreated(t *testing.T, st *fileSystemStorage) {
	t.Helper()

	if _, err := os.Stat(st.path); err == nil {
		return
	} else if os.IsNotExist(err) {
		t.Fatal("storage folder doesn't exist")
	} else {
		log.Fatal("unable to check storage folder existence")
	}
}

func assertBoxFolderIsCreated(t *testing.T, st *fileSystemStorage, bid string) {
	t.Helper()

	if _, err := os.Stat(filepath.Join(st.path, bid)); err == nil {
		return
	} else if os.IsNotExist(err) {
		t.Fatal("box folder doesn't exist")
	} else {
		log.Fatal("unable to check box folder existence")
	}
}

func newCleanUpFileSystemStorageFunc(path string) func() {
	return func() {
		if err := os.RemoveAll(path); err != nil {
			log.Fatal("unable to remove residual data")
		}
	}
}

func TestNewFileSystemStorage(t *testing.T) {
	dir := "tests_box-storage"
	t.Run("create folder and return storage", func(t *testing.T) {
		got := NewFileSystemStorage("", dir)

		assert.NotNil(t, got)
		assertStorageFolderIsCreated(t, got)
	})

	t.Run("panic because empty dirname", func(t *testing.T) {
		assert.Panics(t, func() {
			NewFileSystemStorage("", "")
		})
	})

	t.Cleanup(newCleanUpFileSystemStorageFunc(dir))
}

func TestFileSystemStorage_CreatingBox(t *testing.T) {
	path := ""
	dir := "tests_box-storage"
	st := &fileSystemStorage{&defaulFileSystemHandler{}, filepath.Join(path, dir)}
	createStorageFolder(st)

	t.Run("create box in storage", func(t *testing.T) {
		err := st.CreateBox("box_1")

		assert.Nil(t, err)
		assertBoxFolderIsCreated(t, st, "box_1")
	})

	createBoxFolder(st, "box_2")

	t.Run("returns error because id already exists", func(t *testing.T) {

		err := st.CreateBox("box_2")

		assert.Error(t, err, ErrRepeatedBoxIdentifier)
	})

	t.Run("returns error because unexpected/internal error", func(t *testing.T) {
		st.handler = &mockFileSystemHandler{
			ExistsFunc: func(file string) (bool, error) {
				return false, errFoo
			},
		}
		err := st.CreateBox("box_3")

		assert.Error(t, err, ErrUnableToCreateBox)
	})
	t.Cleanup(newCleanUpFileSystemStorageFunc(st.path))
}

type mockFileSystemHandler struct {
	ExistsFunc func(file string) (bool, error)
	MkdirFunc  func(dirname string) error
}

func (h *mockFileSystemHandler) Exists(file string) (bool, error) {
	return h.ExistsFunc(file)
}

func (h *mockFileSystemHandler) Mkdir(dirname string) error {
	return h.MkdirFunc(dirname)
}
