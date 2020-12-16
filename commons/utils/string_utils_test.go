package utils

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	s, _ := Append("", "")
	fmt.Println(s)

}
