package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	slog "github.com/m2c/kiplestar/commons/log"
	"math/rand"
	"reflect"
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

var special = `-+_!@#$%^&*.,?`

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
	//random special character
	preciousCharacter := string(special[rand.Intn(len(special)-1)])
	result = append(result, string(special[rand.Intn(len(special)-1)]))
	for i := 4; i < length; i++ {
		charValue := GetNextPasswordChar(preciousCharacter)
		result = append(result, charValue)
		preciousCharacter = charValue

	}
	return strings.Join(result, "")
}
func GetNextPasswordChar(previouseChar string) (currentChar string) {
	t := rand.Intn(62)
	if t < 10 {
		currentChar = IndexToChar(previouseChar, 10, func(randIndex int) string { return strconv.Itoa(randIndex) })
	} else if t < 26 {
		currentChar = IndexToChar(previouseChar, 26, func(randIndex int) string { return string(randIndex + 65) })
	} else if t < 46 {
		currentChar = IndexToChar(previouseChar, 26, func(randIndex int) string { return string(randIndex + 97) })
	} else {
		currentChar = IndexToChar(previouseChar, len(special)-1, func(randIndex int) string { return string(special[randIndex]) })
	}
	return currentChar
}
func IndexToChar(previous string, rangeCont int, GetChar func(index int) string) string {
	randCharIndex := rand.Intn(rangeCont)
	currentChar := GetChar(randCharIndex)
	if currentChar == previous { // can not use repetitive characters
		randCharIndex = (randCharIndex + 1) % rangeCont
		currentChar = GetChar(randCharIndex)
	}
	return currentChar
}

var sensitiveWords = []string{
	"password",     //gkuser
	"pin",          //gkuser
	"mobile",       //gkuser
	"phone",        //gkuser
	"account",      //gkuser
	"securityCode", //gkcc
	"number",       //gkcc
}

func SensitiveStruct(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		slog.Errorf("============error to SensitiveStruct:%v", v)
		return ""
	}
	return SensitiveFilter(string(bytes))
}

func containsSensitiveWords(k string) bool {
	for _, kw := range sensitiveWords {
		if strings.Contains(k, kw) {
			return true
		}
	}
	return false
}

func traversalFind(root map[string]interface{}) bool {
	var sensitive bool
	for k, v := range root {
		//Currently, only Map is supported ,Arrays are not currently supported
		if v != nil && reflect.TypeOf(v).Kind() == reflect.Map && traversalFind(v.(map[string]interface{})) {
			sensitive = true
		} else if containsSensitiveWords(k) {
			//Determine the type to avoid errors
			if v != nil && reflect.TypeOf(root[k]).Kind() == reflect.String {
				content := root[k].(string)
				if content != "" {

					// mobile
					if strings.Contains(k, "mobile") ||
						strings.Contains(k, "phone") ||
						strings.Contains(k, "account") {
						if len(content) > 8 {
							root[k] = content[0:2] + "****" + content[len(content)-4:len(content)]
						} else {
							root[k] = "**********"
						}
					} else {
						// other
						root[k] = "**********"
					}
				}
			}
			sensitive = true
		}
	}
	return sensitive
}

func SensitiveFilter(content string) string {
	mapData := make(map[string]interface{})
	err := json.Unmarshal([]byte(content), &mapData)
	if err == nil {
		if traversalFind(mapData) {
			dataByte, err := json.Marshal(mapData)
			if err == nil {
				return string(dataByte)
			}
		}
	}
	return content
}

var isTokenTable = [127]bool{
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'0':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	'A':  true,
	'B':  true,
	'C':  true,
	'D':  true,
	'E':  true,
	'F':  true,
	'G':  true,
	'H':  true,
	'I':  true,
	'J':  true,
	'K':  true,
	'L':  true,
	'M':  true,
	'N':  true,
	'O':  true,
	'P':  true,
	'Q':  true,
	'R':  true,
	'S':  true,
	'T':  true,
	'U':  true,
	'W':  true,
	'V':  true,
	'X':  true,
	'Y':  true,
	'Z':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'a':  true,
	'b':  true,
	'c':  true,
	'd':  true,
	'e':  true,
	'f':  true,
	'g':  true,
	'h':  true,
	'i':  true,
	'j':  true,
	'k':  true,
	'l':  true,
	'm':  true,
	'n':  true,
	'o':  true,
	'p':  true,
	'q':  true,
	'r':  true,
	's':  true,
	't':  true,
	'u':  true,
	'v':  true,
	'w':  true,
	'x':  true,
	'y':  true,
	'z':  true,
	'|':  true,
	'~':  true,
}
var commonHeader map[string]string

const toLower = 'a' - 'A'

func HeaderKey(s string) string {
	upper := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !validHeaderFieldByte(c) {
			return s
		}
		if upper && 'a' <= c && c <= 'z' {
			return canonicalMIMEHeaderKey([]byte(s))
		}
		if !upper && 'A' <= c && c <= 'Z' {
			return canonicalMIMEHeaderKey([]byte(s))
		}
		upper = c == '-'
	}
	return s
}

func validHeaderFieldByte(b byte) bool {
	return int(b) < len(isTokenTable) && isTokenTable[b]
}
func canonicalMIMEHeaderKey(a []byte) string {
	// See if a looks like a header key. If not, return it unchanged.
	for _, c := range a {
		if validHeaderFieldByte(c) {
			continue
		}
		// Don't canonicalize.
		return string(a)
	}

	upper := true
	for i, c := range a {
		// Canonicalize: first letter upper case
		// and upper case after each dash.
		// (Host, User-Agent, If-Modified-Since).
		// MIME headers are ASCII only, so no Unicode issues.
		if upper && 'a' <= c && c <= 'z' {
			c -= toLower
		} else if !upper && 'A' <= c && c <= 'Z' {
			c += toLower
		}
		a[i] = c
		upper = c == '-' // for next time
	}
	// The compiler recognizes m[string(byteSlice)] as a special
	// case, so a copy of a's bytes into a new string does not
	// happen in this map lookup:
	if v := commonHeader[string(a)]; v != "" {
		return v
	}
	return string(a)
}
