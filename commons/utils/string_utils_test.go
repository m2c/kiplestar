package utils

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	s, _ := Append("", "")
	fmt.Println(s)

}

func TestRandomSixString(t *testing.T) {
	fmt.Println(RandomSixString(6))

}
