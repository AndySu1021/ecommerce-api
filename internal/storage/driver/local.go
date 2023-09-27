package driver

import (
	"context"
	"io"
	"os"
	"strings"
)

type Local struct{}

func NewLocal() Client {
	return &Local{}
}

func (d *Local) Upload(ctx context.Context, reader io.Reader, path, filename string) (string, error) {
	filename = getFileName(filename)
	rootDir := "./" + strings.Trim(path, "/") + "/"

	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		return "", err
	}

	out, err := os.Create(rootDir + filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(out, reader)
	if err != nil {
		return "", err
	}

	return getUrl(path, filename), nil
}
