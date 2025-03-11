package usb_test

import (
	"os"
)

type MockFileSystem struct {
	ReadDirFunc      func(name string) ([]os.DirEntry, error)
	EvalSymlinksFunc func(path string) (string, error)
}

func (m MockFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return m.ReadDirFunc(name)
}

func (m MockFileSystem) EvalSymlinks(path string) (string, error) {
	return m.EvalSymlinksFunc(path)
}
