package fileutil

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// at most sniffLen bytes of a file is required to make a decision on its content type
const sniffLen = 512

var imageContentTypes = []string{"image/jpeg", "image/png"}

func IsImageFile(file *multipart.FileHeader) (bool, error) {
	f, err := file.Open()
	if err != nil {
		return false, err
	}
	defer f.Close()

	buff := make([]byte, sniffLen)
	if _, err := f.Read(buff); err != nil {
		return false, err
	}

	return IsImageBuffer(buff), nil
}

func IsImageBuffer(buf []byte) bool {
	contentType := http.DetectContentType(buf)
	return lo.Contains(imageContentTypes, contentType)
}

func ReadBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	if fileHeader == nil {
		return nil, nil
	}

	var byteArray []byte
	file, err := fileHeader.Open()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	byteArray, err = io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}
	_ = file.Close()
	return byteArray, nil
}
