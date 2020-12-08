package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/m2c/kiplestar/commons"
	slog "github.com/m2c/kiplestar/commons/log"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

/**
 * convert to string
 */
func ParseString(b int) string {
	id := strconv.Itoa(b)
	return id
}

func StringToInt64(a string, defaultVal int64) int64 {
	res, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		return defaultVal
	}
	return res
}
func StringToInt(str string, defaultVal int) int {
	res, err := strconv.Atoi(str)
	if err != nil {
		return defaultVal
	}
	return res
}

//get uuid
func GetUuid() string {
	u1 := uuid.NewV4()

	return strings.Replace(u1.String(), "-", "", -1)
}

func StringToTime(timeString string) time.Time {
	t, err := time.ParseInLocation(TIME_LAYOUT, timeString, SetLocation())
	if err != nil {
		slog.Error(err.Error())
	}
	return t
}

func SetLocation() *time.Location {
	local, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		return local
	}
	slog.Errorf(" timestamp err:%s", err.Error())
	return time.Local
}

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//string to md5
func StringToMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func StringToHmac256(secret, s string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func RetryFunction(c func() bool, times int) bool {
	for i := times + 1; i > 0; i-- {
		if c() == true {
			return true
		}
	}
	return false
}
func PushMapNotNull(keyMap map[string]interface{}, key string, value interface{}) {
	if reflect.ValueOf(value).IsZero() == false {
		keyMap[key] = value
	}
}

/**
 * uniform validate parameters and parse the params
 */
func ValidateAndBindParameters(entity interface{}, ctx *iris.Context, info string) (commons.ResponseCode, string) {
	if err := (*ctx).UnmarshalBody(entity, iris.UnmarshalerFunc(json.Unmarshal)); err != nil {
		slog.Errorf("%s error %s", info, err.Error())
		return commons.ParameterError, err.Error()
	}
	if err := Validate(entity); err != nil {
		slog.Errorf("%s error %s", info, err.Error())
		return commons.ValidateError, err.Error()
	}
	return commons.OK, ""
}

/**
*	call the commons service method
 */
func CommonsService(parms BaseParams, f func(parms BaseParams) (interface{}, error)) (interface{}, error) {
	data, err := f(parms)
	slog.Errorf(parms.info+"%s", err.Error())
	return data, err
}

/**
*	base parameters
 */
type BaseParams struct {
	info string
}

func ToString(from interface{}) string {
	v := reflect.ValueOf(from)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	}
	return ""
}
