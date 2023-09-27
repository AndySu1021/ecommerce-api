package helper

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func SaltEncrypt(plaintext, salt string) string {
	if salt == "" {
		salt = "my-salt"
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(salt+plaintext)))
}

func GetUrl(baseUrl, path string) string {
	builder := strings.Builder{}
	builder.WriteString(strings.TrimRight(baseUrl, "/"))
	builder.WriteString("/")
	builder.WriteString(strings.Trim(path, "/"))
	return builder.String()
}

func GenOrderNo(serverId int64) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := time.Now().Unix() / 10
	timestamp := strconv.FormatInt(t, 10)
	timestamp = fmt.Sprintf("%s%03d", timestamp, serverId)
	randNum := strconv.FormatInt(100+r.Int63n(900), 10)
	orderNo := timestamp + randNum
	return orderNo
}

func ConvertToInt64(value interface{}) (int64, error) {
	switch data := value.(type) {
	case []uint8:
		return strconv.ParseInt(string(data), 10, 64)
	case int64:
		return data, nil
	}

	return 0, errors.New("wrong type")
}
