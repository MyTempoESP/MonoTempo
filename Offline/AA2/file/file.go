package file

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

type File struct {
	mu     sync.Mutex
	file   *os.File
	writer *bufio.Writer

	reportChannel chan error

	Caminho string
}

func NewFile(nome string) (a File, err error) {

	a.Caminho = fmt.Sprintf("/tmp/%s", nome)

	f, err := os.OpenFile(a.Caminho, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)

	if err != nil {

		return
	}

	a.file = f

	w := bufio.NewWriter(f)
	a.writer = w

	return
}

func (a *File) Clear() (err error) {

	a.mu.Lock()
	defer a.mu.Unlock()

	err = a.file.Truncate(0)

	if err != nil {

		return
	}

	_, err = a.file.Seek(0, 0)

	return
}

func (a *File) Insert(content string) (err error) {

	a.mu.Lock()
	defer a.mu.Unlock()

	_, err = a.writer.WriteString(fmt.Sprint(content + "\n"))

	return
}

func (a *File) Upload(dest string /* placeholder */) (err error) {

	a.mu.Lock()
	defer a.mu.Unlock()

	err = a.writer.Flush()

	if err != nil {

		return
	}

	err = copyFile(a.Caminho, dest)

	if err != nil {

		log.Println("COPY | Error")
		return
	}

	err = a.file.Truncate(0)

	if err != nil {

		return
	}

	_, err = a.file.Seek(0, 0)

	return
}
