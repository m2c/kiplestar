package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func DataToExcelByte(data interface{}) (rsp []byte, err error) {
	// slice
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() != reflect.Slice {
		return rsp, errors.New("the data must be a slice")
	}

	f := excelize.NewFile()

	//skip ignore fields
	skipRow :=0
	for i := 0; i < t.Elem().NumField(); i++ {
		column := string(rune((i-skipRow) + 65))
		key := fmt.Sprintf("%s%d", column, 1)
		field := t.Elem().Field(i).Name
		if t.Elem().Field(i).Tag.Get("export") == "ignore"{
			skipRow++
			continue
		}

		f.SetCellValue("Sheet1", key, field)

		for iv := 0; iv < v.Len(); iv++ {
			key := fmt.Sprintf("%s%d", column, iv+2)
			f.SetCellValue("Sheet1", key, v.Index(iv).FieldByName(field))
		}

	}

	buff, err := f.WriteToBuffer()
	if err != nil {
		return rsp, err
	}

	return buff.Bytes(), nil
}
