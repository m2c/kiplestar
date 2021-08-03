package utils

import (
	"fmt"
	"testing"
)

func TestStringToMd5(t *testing.T) {
	fmt.Println(StringToMd5("kiplebiz" + "WXnZ7FvZzP"))
	fmt.Println(StringToMd5("kiplepark" + "zCBwKSC2vh"))
	fmt.Println(StringToMd5("common" + "SXz9aMqErL"))
	fmt.Println(StringToMd5("admin" + "NODFIYHlnT"))
}

func TestYuanToMicrometer(t *testing.T) {
	amount := 10000.00
	fmt.Println(YuanToMicrometer(amount))
}
