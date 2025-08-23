package fileutil

import (
	"bytes"
	"mime/multipart"
	"os"
	"path"
)

type File struct {
	Filename string
	Bytes    []byte
	Size     int64
}

// NewFileFromFileHeader creates a new File from multipart.FileHeader.
// If the header is nil it would return nil NOT error,
// but if the content of the header could not be read it would return an error.
func NewFileFromFileHeader(header *multipart.FileHeader) (*File, error) {
	if header == nil {
		return nil, nil
	}

	data, err := ReadBytes(header)
	if err != nil {
		return nil, err
	}

	return &File{
		Filename: header.Filename,
		Bytes:    data,
		Size:     header.Size,
	}, nil
}

func NewFileFromPath(p string) (*File, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	fileName := path.Base(p)

	return &File{
		Filename: fileName,
		Bytes:    data,
		Size:     int64(len(data)),
	}, nil
}

func (f *File) AsBuffer() *bytes.Reader {
	return bytes.NewReader(f.Bytes)
}

func (f *File) IsImage() bool {
	return IsImageBuffer(f.Bytes)
}
