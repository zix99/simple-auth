package box

import (
	"io"
	"os"
)

var Global = NewBox()

func Read(filename string) (io.ReadSeeker, error) {
	return Global.Read(filename)
}

func Stat(filename string) (os.FileInfo, error) {
	return Global.Stat(filename)
}
