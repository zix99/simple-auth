package box

import (
	"io"
)

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

type readSeekCloser struct {
	io.ReadSeeker
}

func (readSeekCloser) Close() error {
	return nil
}
