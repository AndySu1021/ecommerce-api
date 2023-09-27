package driver

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io"
	"strings"
)

type GCS struct {
	Bucket         string `mapstructure:"bucket"`
	CredentialPath string `mapstructure:"credential_path"`
}

func NewGCS(bucket, credentialPath string) Client {
	return &GCS{
		Bucket:         bucket,
		CredentialPath: credentialPath,
	}
}

func (d *GCS) Upload(ctx context.Context, reader io.Reader, path, filename string) (string, error) {
	filename = getFileName(filename)
	rootDir := strings.Trim(path, "/") + "/"

	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(d.CredentialPath))
	if err != nil {
		return "", err
	}

	sw := storageClient.Bucket(d.Bucket).Object(rootDir + filename).NewWriter(ctx)
	if _, err = io.Copy(sw, reader); err != nil {
		return "", err
	}

	if err = sw.Close(); err != nil {
		return "", err
	}

	return getUrl(path, filename), nil
}
