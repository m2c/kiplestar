package httptool

import (
	"bytes"
	"fmt"
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

func FormatRequestItem(key string, tv reflect.Value) (req []RequestItem, err error) {
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
			tmps, err := FormatRequestItem(key, tv.Index(i))
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
