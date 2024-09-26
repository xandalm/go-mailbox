package mailbox

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func createStorage(path, dir string) *fileSystemStorage {
	path = filepath.Join(path, dir)
	st := &fileSystemStorage{
		boxesInfo: []*boxInfo{},
		handler:   &defaulFileSystemHandler{},
	}
	err := os.MkdirAll(path, 0666)
	if err != nil {
		log.Fatal("unable to create storage folder")
		return nil
	}
	st.f, _ = os.Open(path)
	return st
}

func createBox(st *fileSystemStorage, bid string) {
	path := filepath.Join(st.f.Name(), bid)
	if err := os.MkdirAll(path, 0666); err != nil {
		log.Fatalf("unable to create box folder, %v", err)
	}
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("unable to open box file, %v", err)
	}
	pos, _ := slices.BinarySearchFunc(st.boxesInfo, bid, func(e *boxInfo, id string) int {
		return strings.Compare(e.id, id)
	})
	st.boxesInfo = slices.Insert(st.boxesInfo, pos, &boxInfo{bid, f})
}

func createContentFile(st *fileSystemStorage, bid string, cid string, content []byte) {
	f, err := os.OpenFile(filepath.Join(st.f.Name(), bid, cid), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("unable to create content file, %v", err)
	}
	defer f.Close()
	f.Write(content)
}

func assertStorageFolderIsCreated(t *testing.T, st *fileSystemStorage) {
	t.Helper()

	if _, err := os.Stat(st.f.Name()); err == nil {
		return
	} else if os.IsNotExist(err) {
		t.Fatal("storage folder doesn't exist")
	} else {
		log.Fatal("unable to check storage folder existence")
	}
}

func exists(name, msg string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	log.Fatal(msg, err)
	return false
}

func isBoxFolderCreated(st *fileSystemStorage, bid string) bool {
	return exists(filepath.Join(st.f.Name(), bid), "unable to check box folder existence")
}

func isContentFileCreated(st *fileSystemStorage, bid, cid string) bool {
	return exists(filepath.Join(st.f.Name(), bid, cid), "unable to check content file existence")
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

func assertBoxFolderIsEmpty(t *testing.T, st *fileSystemStorage, bid string) {
	t.Helper()

	f, err := os.Open(filepath.Join(st.f.Name(), bid))
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatal("box folder doesn't even exist")
		}
		log.Fatalf("unable to open file, %v", err)
		return
	}
	defer f.Close()
	if _, err := f.Readdirnames(1); err == nil {
		t.Fatal("box folder isn't empty")
	} else if err != io.EOF {
		log.Fatalf("unable to list names inside dir, %v", err)
	}
}

func assertContentFileIsCreated(t *testing.T, st *fileSystemStorage, bid, cid string) {
	t.Helper()

	if !isContentFileCreated(st, bid, cid) {
		t.Fatal("didn't create content file")
	}
}

func assertContentFileDataIsEqual(t *testing.T, st *fileSystemStorage, bid, cid string, want []byte) {
	t.Helper()

	f, err := os.Open(filepath.Join(st.f.Name(), bid, cid))
	if err != nil {
		log.Fatalf("unable to open file, %v", err)
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("unable to read data from file, %v", err)
	}

	if !bytes.Equal(data, want) {
		t.Fatalf("content file data is %v, but want %v", data, want)
	}
}

func newCleanUpStorageFunc(st *fileSystemStorage) func() {
	return func() {
		if st.f == nil {
			return
		}
		for _, box := range st.boxesInfo {
			box.f.Close()
		}
		if err := st.f.Close(); err != nil {
			log.Fatal("unable to close storage file")
		}
	}
}

func TestNewFileSystemStorage(t *testing.T) {
	path := t.TempDir()
	dir := "tests_box-storage"
	t.Run("create folder and return storage", func(t *testing.T) {
		got := NewFileSystemStorage(path, dir)

		assert.NotNil(t, got)
		assert.NotNil(t, got.f)
		assert.Equal(t, got.f.Name(), filepath.Join(path, dir))
		assert.NotNil(t, got.boxesInfo)

		assertStorageFolderIsCreated(t, got)

		t.Cleanup(newCleanUpStorageFunc(got))
	})

	t.Run("panic because empty dirname", func(t *testing.T) {
		assert.Panics(t, func() {
			NewFileSystemStorage("", "")
		})
	})

}

var testDir = "tests_box-storage"

func TestFileSystemStorage_CreatingBox(t *testing.T) {
	path := t.TempDir()
	st := createStorage(path, testDir)

	t.Run("create box in storage", func(t *testing.T) {
		err := st.CreateBox("box_1")

		assert.Nil(t, err)
		_, has := slices.BinarySearchFunc(st.boxesInfo, "box_1", func(e *boxInfo, id string) int {
			return strings.Compare(e.id, id)
		})
		if !has {
			t.Fatal("didn't create box in storage")
		}
		assertBoxFolderIsCreated(t, st, "box_1")
	})

	createBox(st, "box_2")

	t.Run("returns error because id already exists", func(t *testing.T) {

		err := st.CreateBox("box_2")

		assert.Error(t, err, ErrRepeatedBoxIdentifier)
	})

	t.Run("returns error because unexpected/internal error", func(t *testing.T) {
		bid := "box_3"
		t.Run("unable to make box folder", func(t *testing.T) {
			st.handler = &mockFileSystemHandler{
				OpenFunc: dummyOpenFunc,
				ExistsFunc: func(name string) (bool, error) {
					return false, nil
				},
				MkdirFunc: func(name string) error {
					return errFoo
				},
			}
			err := st.CreateBox(bid)

			assert.Error(t, err, ErrUnableToCreateBox)
		})
		t.Run("unable to keep box file open", func(t *testing.T) {
			st.handler = &mockFileSystemHandler{
				OpenFunc: func(name string) (*os.File, error) {
					return nil, errFoo
				},
				ExistsFunc: func(name string) (bool, error) {
					return false, nil
				},
				MkdirFunc: func(name string) error {
					return nil
				},
			}
			err := st.CreateBox(bid)

			assert.Error(t, err, ErrUnableToCreateBox)
		})
	})

	t.Cleanup(newCleanUpStorageFunc(st))
}

func TestFileSystemStorage_ListingBox(t *testing.T) {
	path := t.TempDir()
	st := createStorage(path, testDir)
	createBox(st, "box_A")
	createBox(st, "box_B")
	createBox(st, "box_C")

	t.Run("returns list with id of all boxes in storage", func(t *testing.T) {
		got, err := st.ListBoxes()

		assert.Nil(t, err)
		assert.Contains(t, got, "box_A")
		assert.Contains(t, got, "box_B")
		assert.Contains(t, got, "box_C")
	})

	t.Cleanup(newCleanUpStorageFunc(st))
}

func TestFileSystemStorage_DeletingBox(t *testing.T) {
	path := t.TempDir()
	st := createStorage(path, testDir)
	createBox(st, "box_A")
	createBox(st, "box_B")
	createBox(st, "box_C")
	createContentFile(st, "box_B", "data_1", []byte("foo"))

	t.Run("delete box from storage", func(t *testing.T) {
		err := st.DeleteBox("box_A")

		assert.Nil(t, err)
		assertBoxFolderWasDeleted(t, st, "box_A")
	})

	t.Run("delete box and its contents", func(t *testing.T) {
		err := st.DeleteBox("box_B")

		assert.Nil(t, err)
		assertBoxFolderWasDeleted(t, st, "box_B")
	})

	t.Run("returns error because unexpected/internal error", func(t *testing.T) {
		st.handler = &mockFileSystemHandler{
			OpenFunc: dummyOpenFunc,
			RemoveAllFunc: func(name string) error {
				return errFoo
			},
		}

		err := st.DeleteBox("box_C")

		assert.Error(t, err, ErrUnableToDeleteBox)
	})

	t.Cleanup(newCleanUpStorageFunc(st))
}

func TestFileSystemStorage_CleaningBox(t *testing.T) {
	path := t.TempDir()
	st := createStorage(path, testDir)
	createBox(st, "box_1")
	createContentFile(st, "box_1", "data_1", []byte("foo"))
	createContentFile(st, "box_1", "data_2", []byte("bar"))

	t.Run("clean box data", func(t *testing.T) {
		err := st.CleanBox("box_1")

		assert.Nil(t, err)
		assertBoxFolderIsEmpty(t, st, "box_1")
	})

	t.Run("returns error because unexpected/internal error", func(t *testing.T) {
		st.handler = &mockFileSystemHandler{
			OpenFunc: dummyOpenFunc,
			CleanFunc: func(dir *os.File) error {
				return errFoo
			},
		}

		err := st.CleanBox("box_1")

		assert.Error(t, err, ErrUnableToCleanBox)
	})

	t.Cleanup(newCleanUpStorageFunc(st))
}

func TestFileSystemStorage_CreatingContent(t *testing.T) {
	path := t.TempDir()
	st := createStorage(path, testDir)
	createBox(st, "box_1")

	t.Run("create content", func(t *testing.T) {
		bid := "box_1"
		cid := "data_1"
		data := []byte("foo")
		err := st.CreateContent(bid, cid, data)

		assert.Nil(t, err)
		assertContentFileIsCreated(t, st, bid, cid)
		assertContentFileDataIsEqual(t, st, bid, cid, data)
	})

	createContentFile(st, "box_1", "data_2", []byte("bar"))

	t.Run("returns error because box doesn't exist", func(t *testing.T) {
		err := st.CreateContent("box_2", "data_2", []byte("baz"))

		assert.Error(t, err, ErrBoxNotFoundToPostContent)
	})

	t.Run("returns error because id already exists", func(t *testing.T) {
		err := st.CreateContent("box_1", "data_2", []byte("baz"))

		assert.Error(t, err, ErrRepeatedContentIdentifier)
	})

	t.Run("returns error because unexpected/internal error", func(t *testing.T) {
		bid := "box_1"
		cid := "data_3"
		data := []byte("baz")
		t.Run("unable to check id existence", func(t *testing.T) {
			st.handler = &mockFileSystemHandler{
				OpenFunc: dummyOpenFunc,
				ExistsFunc: func(name string) (bool, error) {
					return false, errFoo
				},
			}
			err := st.CreateContent(bid, cid, data)

			assert.Error(t, err, ErrUnableToPostContent)
		})
		t.Run("unable to write file", func(t *testing.T) {
			st.handler = &mockFileSystemHandler{
				OpenFunc: dummyOpenFunc,
				ExistsFunc: func(name string) (bool, error) {
					return false, nil
				},
				WriteFileFunc: func(name string, data []byte) error {
					return errFoo
				},
			}
			err := st.CreateContent(bid, cid, data)

			assert.Error(t, err, ErrUnableToPostContent)
		})
	})

	t.Cleanup(newCleanUpStorageFunc(st))
}

func TestFileSystemStorage_ReadingContent(t *testing.T) {
	path := t.TempDir()
	st := createStorage(path, testDir)
	createBox(st, "box_1")
	createContentFile(st, "box_1", "data_1", []byte("foo"))

	t.Run("read content", func(t *testing.T) {
		got, err := st.ReadContent("box_1", "data_1")

		assert.Nil(t, err)
		assert.NotNil(t, got)
		want := []byte("foo")
		assert.Equal(t, got, want)
	})

	t.Run("returns error because box doesn't exist", func(t *testing.T) {
		data, err := st.ReadContent("box_2", "data_2")

		assert.Nil(t, data)
		assert.Error(t, err, ErrBoxNotFoundToReadContent)
	})

	t.Run("returns error because content file doesn't exist", func(t *testing.T) {
		data, err := st.ReadContent("box_1", "data_2")

		assert.Nil(t, data)
		assert.Error(t, err, ErrContentNotFound)
	})

	t.Run("returns error because unexpected/internal error", func(t *testing.T) {
		bid := "box_1"
		cid := "data_1"
		t.Run("unable to check content existence", func(t *testing.T) {
			st.handler = &mockFileSystemHandler{
				OpenFunc: dummyOpenFunc,
				ExistsFunc: func(name string) (bool, error) {
					return false, errFoo
				},
			}
			data, err := st.ReadContent(bid, cid)

			assert.Nil(t, data)
			assert.Error(t, err, ErrUnableToReadContent)
		})
		t.Run("unable to read file", func(t *testing.T) {
			st.handler = &mockFileSystemHandler{
				OpenFunc: dummyOpenFunc,
				ExistsFunc: func(name string) (bool, error) {
					return true, nil
				},
				ReadFileFunc: func(name string) ([]byte, error) {
					return nil, errFoo
				},
			}
			data, err := st.ReadContent(bid, cid)

			assert.Nil(t, data)
			assert.Error(t, err, ErrUnableToReadContent)
		})
	})

	t.Cleanup(newCleanUpStorageFunc(st))
}

var dummyOpenFunc = func(_ string) (*os.File, error) {
	return nil, nil
}

type mockFileSystemHandler struct {
	OpenFunc      func(name string) (*os.File, error)
	ExistsFunc    func(name string) (bool, error)
	MkdirFunc     func(name string) error
	RemoveFunc    func(name string) error
	RemoveAllFunc func(name string) error
	CleanFunc     func(dir *os.File) error
	WriteFileFunc func(name string, data []byte) error
	ReadFileFunc  func(name string) ([]byte, error)
}

func (h *mockFileSystemHandler) Open(name string) (*os.File, error) {
	return h.OpenFunc(name)
}

func (h *mockFileSystemHandler) Exists(name string) (bool, error) {
	return h.ExistsFunc(name)
}

func (h *mockFileSystemHandler) Mkdir(name string) error {
	return h.MkdirFunc(name)
}

func (h *mockFileSystemHandler) Remove(name string) error {
	return h.RemoveFunc(name)
}

func (h *mockFileSystemHandler) RemoveAll(name string) error {
	return h.RemoveAllFunc(name)
}

func (h *mockFileSystemHandler) Clean(dir *os.File) error {
	return h.CleanFunc(dir)
}

func (h *mockFileSystemHandler) WriteFile(name string, data []byte) error {
	return h.WriteFileFunc(name, data)
}

func (h *mockFileSystemHandler) ReadFile(name string) ([]byte, error) {
	return h.ReadFileFunc(name)
}
