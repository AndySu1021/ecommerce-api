package helper

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func ParseEcPayRequest(body string, params interface{}) error {
	decodedValue, err := url.QueryUnescape(body)
	if err != nil {
		return err
	}

	parseResult, err := url.ParseQuery(decodedValue)
	if err != nil {
		return err
	}

	dType := reflect.TypeOf(params)
	dhVal := reflect.ValueOf(params)

	for i := 0; i < dType.Elem().NumField(); i++ {
		field := dType.Elem().Field(i)
		key := field.Tag.Get("json")
		val := parseResult.Get(key)
		kind := field.Type.Kind()

		switch kind {
		case reflect.String:
			dhVal.Elem().Field(i).SetString(val)
		case reflect.Int64:
			tmpVal, convErr := strconv.ParseInt(val, 10, 64)
			if convErr != nil {
				return convErr
			}
			dhVal.Elem().Field(i).SetInt(tmpVal)
		case reflect.Int:
			tmpVal, convErr := strconv.ParseInt(val, 10, 64)
			if convErr != nil {
				return convErr
			}
			dhVal.Elem().Field(i).SetInt(tmpVal)
		default:
			return fmt.Errorf("wrong type")
		}
	}

	return nil
}

func GetCheckMacValue(params interface{}, key, iv string) (string, error) {
	encodedStr, err := GetEncodedStr(params)
	if err != nil {
		return "", err
	}

	// generate url encoded parameters
	encodedStr = fmt.Sprintf("HashKey=%s&%s&HashIV=%s", key, encodedStr, iv)

	// generate url encoded string
	plaintext := url.QueryEscape(encodedStr)
	plaintext = strings.ReplaceAll(plaintext, "%2d", "-")
	plaintext = strings.ReplaceAll(plaintext, "%5f", "_")
	plaintext = strings.ReplaceAll(plaintext, "%2e", ".")
	plaintext = strings.ReplaceAll(plaintext, "%21", "!")
	plaintext = strings.ReplaceAll(plaintext, "%2A", "*")
	plaintext = strings.ReplaceAll(plaintext, "%28", "(")
	plaintext = strings.ReplaceAll(plaintext, "%29", ")")
	plaintext = strings.ToLower(plaintext)

	// generate md5 hash
	hash := fmt.Sprintf("%x", md5.Sum([]byte(plaintext)))
	hash = strings.ToUpper(hash)

	return hash, nil
}

func GetEncodedStr(params interface{}) (string, error) {
	tmpMap := make(map[string]interface{})
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(jsonStr, &tmpMap); err != nil {
		return "", err
	}

	// sort by ascii order
	keys := make([]string, 0)
	for k := range tmpMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	urlEncodedStr := make([]string, 0)
	for _, key := range keys {
		urlEncodedStr = append(urlEncodedStr, fmt.Sprintf("%s=%v", key, tmpMap[key]))
	}

	return strings.Join(urlEncodedStr, "&"), nil
}
