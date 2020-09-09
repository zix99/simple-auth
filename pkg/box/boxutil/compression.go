package boxutil

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Compress(b []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	gzw := gzip.NewWriter(&buf)

	if _, err := gzw.Write(b); err != nil {
		return nil, err
	}
	if err := gzw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decompress(b []byte) ([]byte, error) {
	gzw, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(gzw)
}

func Must(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}
