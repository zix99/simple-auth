package box

import (
	"os"
)

var Global = NewBox()

func Read(filename string) (ReadSeekCloser, error) {
	return Global.Read(filename)
}

func Stat(filename string) (os.FileInfo, error) {
	return Global.Stat(filename)
}
