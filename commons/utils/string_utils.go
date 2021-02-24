package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

func RandomSixString(length int) string {
	// 48 ~ 57 数字
	// 65 ~ 90 A ~ Z //26
	// 97 ~ 122 a ~ z //26
	// A total of 62 characters, random from 0 to 61, when less than 10, random in the number range, [一共62个字符，在0~61进行随机，小于10时，在数字范围随机，]
	// Less than 36 are random in uppercase range, others are random in lowercase range[小于36在大写范围内随机，其他在小写范围随机]
	rand.Seed(time.Now().UnixNano())
	result := make([]string, 0, length)
	//uppercase
	result = append(result, string(rand.Intn(26)+65))
	//lowercase
	result = append(result, string(rand.Intn(26)+97))
	//random number
	result = append(result, strconv.Itoa(rand.Intn(10)))
	for i := 3; i < length; i++ {
		t := rand.Intn(62)
		if t < 10 {
			result = append(result, strconv.Itoa(rand.Intn(10)))
		} else if t < 36 {
			result = append(result, string(rand.Intn(26)+65))
		} else {
			result = append(result, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(result, "")
}

var sensitiveFields = []string{"password", "confirm_password", "old_password", "pin"}

func SensitiveFilter(content string) string {
	mapData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(content), &mapData); err == nil {
		var sensitive bool
		for i := range sensitiveFields {
			if _, ok := mapData[sensitiveFields[i]]; ok {
				mapData[sensitiveFields[i]] = "**********"
				sensitive = true
			}
		}
		if sensitive {
			if dataByte, err := json.Marshal(mapData); err == nil {
				return string(dataByte)
			}
		}
	}
	return content
}
