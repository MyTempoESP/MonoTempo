package usb

import (
	"os"
	"path/filepath"
)

type FileSystem interface {
	ReadDir(name string) ([]os.DirEntry, error)
	EvalSymlinks(path string) (string, error)
}

type OSFileSystem struct{}

func (OSFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (OSFileSystem) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}
