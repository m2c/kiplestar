package utils

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func DeepFields(interfaceType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < interfaceType.NumField(); i++ {
		v := interfaceType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, DeepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}

	return fields
}
func StructCopy(DstStructPtr interface{}, SrcStructPtr interface{}) {
	srcv := reflect.ValueOf(SrcStructPtr)
	dstv := reflect.ValueOf(DstStructPtr)
	srct := reflect.TypeOf(SrcStructPtr)
	dstt := reflect.TypeOf(DstStructPtr)
	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		panic("Fatal error:type of parameters must be Ptr of value")
	}
	if srcv.IsNil() || dstv.IsNil() {
		panic("Fatal error:value of parameters should not be nil")
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()
	srcfields := DeepFields(reflect.ValueOf(SrcStructPtr).Elem().Type())
	for _, v := range srcfields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.IsZero() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return
}

//transfer url.Values to struct
func Transfer(values url.Values, s interface{}) error {
	val := reflect.ValueOf(s)
	va_pre := val
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.New("Transfer() expects struct input. ")
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("Transfer() expects struct input. ")
	}
	return reflectValueFromTag(values, val, va_pre)
}

//transfer url.Values to struct
func TransferByParam(values url.Values, s interface{}) error {
	val := reflect.ValueOf(s)
	va_pre := val
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.New("Transfer() expects struct input. ")
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("Transfer() expects struct input. ")
	}
	return reflectValue(values, val, va_pre)
}
func reflectValueFromTag(values url.Values, val reflect.Value, va_pre reflect.Value) error {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		kt := typ.Field(i)
		tag := kt.Tag.Get("param")
		if tag == "-" {
			continue
		}
		sv := val.Field(i)
		uv := getVal(values, tag, true)
		switch sv.Kind() {
		case reflect.String:
			sv.SetString(uv)
		case reflect.Bool:
			b, err := strconv.ParseBool(uv)
			if err != nil {
				return errors.New(fmt.Sprintf("cast bool has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
			}
			sv.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(uv, 10, 64)
			if err != nil || sv.OverflowUint(n) {
				return errors.New(fmt.Sprintf("cast uint has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
			}
			sv.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(uv, 10, 64)
			if err != nil || sv.OverflowInt(n) {
				return errors.New(fmt.Sprintf("cast int has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
			}
			sv.SetInt(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(uv, sv.Type().Bits())
			if err != nil || sv.OverflowFloat(n) {
				return errors.New(fmt.Sprintf("cast float has error, expect type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
			}
			sv.SetFloat(n)
		case reflect.Struct:
			if "time.Time" == typ.Field(i).Type.String() {
				t, er := time.Parse("2006-01-02T15:04:05Z", uv)
				if er == nil {
					//va_pre.Elem().FieldByName(kt.Name).Set(reflect.ValueOf(t))
					sv.Set(reflect.ValueOf(t))
				}

			}
		default:
			return errors.New(fmt.Sprintf("unsupported type: %v ,val: %v ,query key: %v", sv.Type(), uv, tag))
		}
	}
	return nil
}
func reflectValue(values url.Values, val reflect.Value, va_pre reflect.Value) error {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		kt := typ.Field(i)
		sv := val.Field(i)
		uv := getVal(values, kt.Name, false)
		switch sv.Kind() {
		case reflect.String:
			sv.SetString(uv)
		case reflect.Bool:
			b, err := strconv.ParseBool(uv)
			if err != nil {
				return errors.New(fmt.Sprintf("cast bool has error, expect type: %v ,val: %v ", sv.Type(), uv))
			}
			sv.SetBool(b)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(uv, 10, 64)
			if err != nil || sv.OverflowUint(n) {
				return errors.New(fmt.Sprintf("cast uint has error, expect type: %v ,val: %v ", sv.Type(), uv))
			}
			sv.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(uv, 10, 64)
			if err != nil || sv.OverflowInt(n) {
				return errors.New(fmt.Sprintf("cast int has error, expect type: %v ,val: %v ", sv.Type(), uv))
			}
			sv.SetInt(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(uv, sv.Type().Bits())
			if err != nil || sv.OverflowFloat(n) {
				return errors.New(fmt.Sprintf("cast float has error, expect type: %v ,val: %v", sv.Type(), uv))
			}
			sv.SetFloat(n)
		case reflect.Struct:
			if "time.Time" == typ.Field(i).Type.String() {
				t, er := time.Parse("2006-01-02T15:04:05Z", uv)
				if er == nil {
					va_pre.Elem().FieldByName(kt.Name).Set(reflect.ValueOf(t))
				}
			}
		default:
			return errors.New(fmt.Sprintf("unsupported type: %v ,val: %v", sv.Type(), uv))
		}
	}
	return nil
}

//get val, if absent get from tag default val
func getVal(values url.Values, tagOrName string, tagType bool) string {

	if tagType {
		name, opts := parseTag(tagOrName)
		uv := values.Get(name)
		optsLen := len(opts)
		if optsLen > 0 {
			if optsLen == 1 && uv == "" {
				uv = opts[0]
			}
		}
		return uv
	} else {
		uv := values.Get(tagOrName)
		return uv
	}
}

type tagOptions []string

func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}
