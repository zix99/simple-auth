package box

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

type EmbedFile struct {
	os.FileInfo

	name    string
	size    int64
	modtime int64

	bytes      []byte
	compressed bool
}

type EmbedBox struct {
	Box
	files     map[string]*EmbedFile
	CheckDisk bool
	Verbose   bool
}

type Box interface {
	Read(filename string) (ReadSeekCloser, error)
	Stat(filename string) (os.FileInfo, error)
}

func NewBox() *EmbedBox {
	return &EmbedBox{
		files:     make(map[string]*EmbedFile),
		CheckDisk: true,
		Verbose:   false,
	}
}

func (s *EmbedFile) Name() string {
	return s.name
}

func (s *EmbedFile) Size() int64 {
	return int64(len(s.bytes))
}

func (s *EmbedFile) Mode() os.FileMode {
	return os.ModePerm
}

func (s *EmbedFile) ModTime() time.Time {
	return time.Unix(s.modtime, 0)
}

func (s *EmbedFile) IsDir() bool {
	return false
}

func (s *EmbedFile) Sys() interface{} {
	return nil
}

func (s *EmbedFile) Read() io.ReadSeeker {
	return bytes.NewReader(s.bytes)
}

func (s *EmbedBox) Add(filename string, file *EmbedFile) {
	filename = path.Clean(filename)
	s.files[filename] = file
	logrus.Debugf("Added file %s (%d bytes) to box", filename, file.size)
}

func (s *EmbedBox) AddBytes(filename string, b []byte) {
	s.files[filename] = &EmbedFile{
		bytes:      b,
		compressed: false,
		modtime:    0,
		name:       path.Base(filename),
		size:       int64(len(b)),
	}
}

func (s *EmbedBox) ReadEx(filename string, alwaysCheckDisk bool) (ReadSeekCloser, error) {
	filename = path.Clean(filename)

	if s.CheckDisk || alwaysCheckDisk {
		if f, err := os.Open(filename); err == nil {
			if s.Verbose {
				logrus.Infof("Opened file %s from disk", filename)
			}
			return f, nil
		}
	}

	if f, ok := s.files[filename]; ok {
		if s.Verbose {
			logrus.Infof("Opened file %s from box", filename)
		}
		return readSeekCloser{f.Read()}, nil
	}

	if s.Verbose {
		logrus.Warnf("File not found %s", filename)
	}

	return nil, errors.New("File not found")
}

func (s *EmbedBox) Read(filename string) (ReadSeekCloser, error) {
	return s.ReadEx(filename, false)
}

func (s *EmbedBox) Stat(filename string) (os.FileInfo, error) {
	filename = path.Clean(filename)

	if s.CheckDisk {
		if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
			return fi, nil
		}
	}

	if f, ok := s.files[filename]; ok {
		return f, nil
	}
	return nil, errors.New("File not found")
}
