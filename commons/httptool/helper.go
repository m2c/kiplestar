package httptool

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

// letter is upper or not
func IsUpperLetter(letter rune) bool {
	if letter >= 'A' && letter <= 'Z' {
		return true
	} else {
		return false
	}
}

func IsLowerLetter(letter rune) bool {
	if letter >= 'a' && letter <= 'z' {
		return true
	} else {
		return false
	}
}

func TransLetterToUpper(letter rune) string {
	if IsLowerLetter(letter) {
		letter -= 'a' - 'A'
	}
	return string(letter)
}

func TransLetterToLower(letter rune) string {
	if IsUpperLetter(letter) {
		letter += 'a' - 'A'
	}
	return string(letter)
}

// like transform "to_lower_snake_case" to "toLowerSnakeCase"
func ToLowerCamelCase(s string) string {
	var dst bytes.Buffer
	var flag bool
	for index, letter := range s {
		if index == 0 {
			dst.WriteString(TransLetterToLower(letter))
		} else if letter == '_' || letter == '-' {
			flag = true
		} else if flag {
			flag = false
			dst.WriteString(TransLetterToUpper(letter))
		} else {
			dst.WriteString(string(letter))
		}
	}

	return dst.String()
}

func getJsonName(f reflect.StructField) string {
	name := f.Tag.Get("json")
	if name == "" {
		name = f.Name
	}
	return name
}

func formatRequestItem(key string, tv reflect.Value) (req []RequestItem, err error) {
	tmp := RequestItem{
		Key: key,
	}
	switch tv.Kind() {
	case reflect.String:
		tmp.Value = tv.String()
		req = []RequestItem{tmp}
	case reflect.Bool:
		tmp.Value = strconv.FormatBool(tv.Bool())
		req = []RequestItem{tmp}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tmp.Value = fmt.Sprintf("%d", tv.Int())
		req = []RequestItem{tmp}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		tmp.Value = fmt.Sprintf("%d", tv.Uint())
		req = []RequestItem{tmp}
	case reflect.Float32, reflect.Float64:
		tmp.Value = fmt.Sprintf("%f", tv.Float())
		req = []RequestItem{tmp}
	case reflect.Slice:
		for i := 0; i < tv.Len(); i++ {
			tmps, err := formatRequestItem(key, tv.Index(i))
			if err != nil {
				return req, err
			}
			req = append(req, tmps...)
		}
		return
	case reflect.Struct:
		switch v := tv.Interface().(type) {
		case time.Time:
			tmp.Value = v.Format(time.RFC3339)
			req = []RequestItem{tmp}
		default:
			err = fmt.Errorf("field [%s] type is invalid.", key)
			return
		}
	default:
		err = fmt.Errorf("field [%s] type is invalid.", key)
		return
	}
	return
}

func FormatRequestParams(req interface{}) ([]RequestItem, error) {
	rs := []RequestItem{}
	tv := reflect.ValueOf(req)
	tp := tv.Type()
	switch tp.Kind() {
	case reflect.Slice:
		for i := 0; i < tp.NumField(); i++ {
			tmps, err := formatRequestItem(getJsonName(tp.Field(i)), tv.Index(i))
			if err != nil {
				return rs, err
			}
			rs = append(rs, tmps...)
		}
	case reflect.Struct:
		for i := 0; i < tp.NumField(); i++ {
			tmps, err := formatRequestItem(getJsonName(tp.Field(i)), tv.Field(i))
			if err != nil {
				return rs, err
			}
			rs = append(rs, tmps...)
		}
	case reflect.Map:
		keys := tv.MapKeys()
		for _, key := range keys {
			tmps, err := formatRequestItem(key.String(), tv.MapIndex(key))
			if err != nil {
				return rs, err
			}
			rs = append(rs, tmps...)
		}
	default:
		return rs, errors.New("req's type is invalid")
	}

	return rs, nil
}

func FormatQueryUrl(urlStr string, req interface{}) (string, error) {
	ul, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	rs, err := FormatRequestParams(req)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	for _, v := range rs {
		params.Add(v.Key, v.Value)
	}
	if ul.RawQuery == "" {
		ul.RawQuery = params.Encode()
	} else {
		ul.RawQuery += "&" + params.Encode()
	}
	return ul.String(), nil
}

func ConvToMap(params interface{}) (req map[string]string, err error) {
	tv := reflect.ValueOf(params)
	if tv.Kind() == reflect.Ptr {
		tv = tv.Elem()
	}
	switch tv.Kind() {
	case reflect.Struct:
		req = map[string]string{}
		tp := tv.Type()
		n := tv.NumField()
		for i := 0; i < n; i++ {
			tmps, err := formatRequestItem(getJsonName(tp.Field(i)), tv.Field(i))
			if err != nil {
				return req, err
			}
			for _, tmp := range tmps {
				req[tmp.Key] = tmp.Value
			}
		}
	case reflect.Map:
		v, ok := params.(map[string]string)
		if !ok {
			err = errors.New("params type is invalid, only struct and map is supported.")
			return
		}
		req = v
	default:
		err = errors.New("params type is invalid, only struct and map is supported.")
		return
	}
	return
}
