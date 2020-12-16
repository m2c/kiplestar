package utils

import (
	"bytes"
	"errors"
)

func Append(source string, strings ...string) (string, error) {
	var buffer bytes.Buffer
	_, err := buffer.WriteString(source)
	if err != nil {

		return "", errors.New("append string has something wrong ")
	}
	for _, value := range strings {
		_, err1 := buffer.WriteString(value)
		if err1 != nil {
			return "", errors.New("append string has something wrong ")
		}
	}
	return buffer.String(), nil
}
