package driver

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

type Client interface {
	Upload(ctx context.Context, reader io.Reader, path, filename string) (string, error)
}

func getFileName(filename string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tmpName := strings.Split(filename, ".")
	ext := tmpName[len(tmpName)-1]
	return fmt.Sprintf("%d.%s", time.Now().UnixMilli()*1000+int64(r.Intn(999)+1), ext)
}

func getUrl(path, filename string) string {
	// url without domain
	builder := strings.Builder{}
	builder.WriteString(strings.Trim(path, "/"))
	builder.WriteString("/")
	builder.WriteString(filename)

	return builder.String()
}
