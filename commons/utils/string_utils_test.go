package utils

import (
	"testing"
)

func TestFindReplaceString(t *testing.T) {
	replaceString, _ := FindReplaceString("我是中国人", []string{"中国人"}, "魅力")
	t.Log(replaceString)
}
