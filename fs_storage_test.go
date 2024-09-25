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

func isBoxFolderCreated(st *fileSystemStorage, bid string) bool {
	if _, err := os.Stat(filepath.Join(st.path, bid)); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	log.Fatal("unable to check box folder existence")
	return false
}

func assertBoxFolderIsCreated(t *testing.T, st *fileSystemStorage, bid string) {
	t.Helper()

	if !isBoxFolderCreated(st, bid) {
		t.Fatal("box folder doesn't exist")
	}
}

func assertBoxFolderWasDeleted(t *testing.T, st *fileSystemStorage, bid string) {
	t.Helper()

	if isBoxFolderCreated(st, bid) {
		t.Fatal("box folder still exists")
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

var testPath = ""
var testDir = "tests_box-storage"

func TestFileSystemStorage_CreatingBox(t *testing.T) {
	st := &fileSystemStorage{&defaulFileSystemHandler{}, filepath.Join(testPath, testDir)}
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
			ExistsFunc: func(name string) (bool, error) {
				return false, errFoo
			},
		}
		err := st.CreateBox("box_3")

		assert.Error(t, err, ErrUnableToCreateBox)
	})

	t.Cleanup(newCleanUpFileSystemStorageFunc(st.path))
}

func TestFileSystemStorage_ListingBox(t *testing.T) {
	st := &fileSystemStorage{&defaulFileSystemHandler{}, filepath.Join(testPath, testDir)}
	createStorageFolder(st)
	createBoxFolder(st, "box_A")
	createBoxFolder(st, "box_B")
	createBoxFolder(st, "box_C")

	t.Run("returns list with id of all boxes in storage", func(t *testing.T) {
		got, err := st.ListBoxes()

		assert.Nil(t, err)
		assert.Contains(t, got, "box_A")
		assert.Contains(t, got, "box_B")
		assert.Contains(t, got, "box_C")
	})

	t.Run("returns error because unexpected/internal error", func(t *testing.T) {
		st.handler = &mockFileSystemHandler{
			LsFunc: func(name string) ([]string, error) {
				return nil, errFoo
			},
		}

		boxes, err := st.ListBoxes()

		assert.Nil(t, boxes)
		assert.Error(t, err, ErrUnableToListBoxes)
	})

	t.Cleanup(newCleanUpFileSystemStorageFunc(st.path))
}

func TestFileSystemStorage_DeletingBox(t *testing.T) {
	st := &fileSystemStorage{&defaulFileSystemHandler{}, filepath.Join(testPath, testDir)}
	createStorageFolder(st)
	createBoxFolder(st, "box_A")
	createBoxFolder(st, "box_B")

	t.Run("delete box from storage", func(t *testing.T) {
		err := st.DeleteBox("box_A")

		assert.Nil(t, err)
		assertBoxFolderWasDeleted(t, st, "box_A")
	})

	t.Run("returns error because unexpected/internal error", func(t *testing.T) {
		st.handler = &mockFileSystemHandler{
			RemoveFunc: func(name string) error {
				return errFoo
			},
		}

		err := st.DeleteBox("box_B")

		assert.Error(t, err, ErrUnableToDeleteBox)
	})

	t.Cleanup(newCleanUpFileSystemStorageFunc(st.path))
}

type mockFileSystemHandler struct {
	ExistsFunc func(name string) (bool, error)
	MkdirFunc  func(name string) error
	LsFunc     func(name string) ([]string, error)
	RemoveFunc func(name string) error
}

func (h *mockFileSystemHandler) Exists(name string) (bool, error) {
	return h.ExistsFunc(name)
}

func (h *mockFileSystemHandler) Mkdir(name string) error {
	return h.MkdirFunc(name)
}

func (h *mockFileSystemHandler) Ls(name string) ([]string, error) {
	return h.LsFunc(name)
}

func (h *mockFileSystemHandler) Remove(name string) error {
	return h.RemoveFunc(name)
}
