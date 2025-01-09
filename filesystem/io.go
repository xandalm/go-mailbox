package filesystem

import (
	"io/fs"
	"os"
)

type ioResult[T any] struct {
	data T
	err  error
}

func getDirectoryNamesFn[T []string](ch chan ioResult[T], f *os.File) {
	f.Seek(0, 0)
	names, err := f.Readdirnames(0)
	ch <- ioResult[T]{names, err}
}

func getDirectoryNames(f *os.File) chan ioResult[[]string] {
	ch := make(chan ioResult[[]string], 1)
	go getDirectoryNamesFn(ch, f)
	return ch
}

func getDirectoryEntriesFn[T []fs.DirEntry](ch chan ioResult[T], f *os.File) {
	f.Seek(0, 0)
	entries, err := f.ReadDir(0)
	ch <- ioResult[T]{entries, err}
}

func getDirectoryEntries(f *os.File) chan ioResult[[]fs.DirEntry] {
	ch := make(chan ioResult[[]fs.DirEntry], 1)
	go getDirectoryEntriesFn(ch, f)
	return ch
}

func getFileInfoFromDirEntryFn[T *fs.FileInfo](ch chan ioResult[T], e fs.DirEntry) {
	info, err := e.Info()
	ch <- ioResult[T]{&info, err}
}

func getFileInfoFromDirEntry(e fs.DirEntry) chan ioResult[*fs.FileInfo] {
	ch := make(chan ioResult[*fs.FileInfo], 1)
	go getFileInfoFromDirEntryFn(ch, e)
	return ch
}

func getFileInfoFn[T *fs.FileInfo](ch chan ioResult[T], name string) {
	info, err := os.Stat(name)
	ch <- ioResult[T]{&info, err}
}

func getFileInfo(name string) chan ioResult[*fs.FileInfo] {
	ch := make(chan ioResult[*fs.FileInfo], 1)
	go getFileInfoFn(ch, name)
	return ch
}

func readFileContentFn[T []byte](ch chan ioResult[T], name string) {
	data, err := os.ReadFile(name)
	ch <- ioResult[T]{data, err}
}

func readFileContent(name string) chan ioResult[[]byte] {
	ch := make(chan ioResult[[]byte], 1)
	go readFileContentFn(ch, name)
	return ch
}

func createDirFn(ch chan error, name string) {
	err := os.Mkdir(name, 0666)
	ch <- err
}

func createDir(name string) chan error {
	ch := make(chan error)
	go createDirFn(ch, name)
	return ch
}

func openFileRWFn[T *os.File](ch chan ioResult[T], name string) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	ch <- ioResult[T]{f, err}
}

func openFileRW(name string) chan ioResult[*os.File] {
	ch := make(chan ioResult[*os.File], 1)
	go openFileRWFn(ch, name)
	return ch
}

func openFileFn[T *os.File](ch chan ioResult[T], name string) {
	f, err := os.Open(name)
	ch <- ioResult[T]{f, err}
}

func openFile(name string) chan ioResult[*os.File] {
	ch := make(chan ioResult[*os.File], 1)
	go openFileFn(ch, name)
	return ch
}

func writeContentFn(ch chan error, f *os.File, c []byte) {
	_, err := f.Write(c)
	ch <- err
}

func writeContent(f *os.File, c []byte) chan error {
	ch := make(chan error, 1)
	go writeContentFn(ch, f, c)
	return ch
}

func removeFn(ch chan error, name string) {
	err := os.RemoveAll(name)
	ch <- err
}

func remove(name string) chan error {
	ch := make(chan error, 1)
	go removeFn(ch, name)
	return ch
}
