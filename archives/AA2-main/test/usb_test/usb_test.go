package usb_test

import (
	"errors"
	"os"
	"testing"

	"aa2/usb"
)

type MockDirEntry struct {
	name string
}

func (m MockDirEntry) Name() string               { return m.name }
func (m MockDirEntry) IsDir() bool                { return false }
func (m MockDirEntry) Type() os.FileMode          { return 0 }
func (m MockDirEntry) Info() (os.FileInfo, error) { return nil, nil }

func TestUSB(t *testing.T) {

	mockFS := MockFileSystem{
		ReadDirFunc: func(name string) ([]os.DirEntry, error) {

			if name == "/sys/block/" {
				return []os.DirEntry{
					MockDirEntry{name: "sda"},
					MockDirEntry{name: "sdb"},
				}, nil
			}

			return nil, errors.New("unexpected path")
		},

		EvalSymlinksFunc: func(path string) (string, error) {

			if path == "/sys/block/sda/device" {
				return "/devices/pci0000:00/usb1", nil
			}

			if path == "/sys/block/sdb/device" {
				return "/devices/pci0000:00/nonusb", nil
			}

			return "", errors.New("unexpected path")
		},
	}

	_, check, err := usb.CheckUSBStorageDevice(mockFS)

	if err != nil {

		t.Fatalf("unexpected error: %v", err)
	}

	if !check {

		t.Fatalf("expected true, got false")
	}
}
