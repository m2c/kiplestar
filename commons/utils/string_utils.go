package utils

import (
	"errors"
	"fmt"
	"strings"
)

/**
 * Sensitive word replacement
 */
func FindReplaceString(s string, words []string, replace string) (string, error) {

	length := len(s)
	if length == 0 {
		return "", errors.New(" the target string is empty")
	}
	word_len := len(words)
	for i := 0; i < word_len; i++ {
		word := words[i]
		word_int_length := []rune(s)
		fmt.Println(len(word_int_length))

		replace_int_word := []rune(replace)
		fmt.Println(len(replace_int_word))

		index := strings.Index(s, word)
		//从第7个开始替换，替换为
		fmt.Println(index, len([]rune(s)))

		//替换的敏感词和替换的字符一样长
		/*if len(word_int_length)==len(replace_int_word)  {
			index:=strings.Index(s,word)
		}*/

		/*

			index:=strings.Index(s,word)
			temp1 := []rune(word)
			fmt.Println(len(temp1))*/

		//length := len(temp)
		s = strings.Replace(s, word, replace, -1)

		/*index := strings.Index(s, word)*/

		/*var myNum []string
		if len(replace)==1 {
			for j:=0;j<len(word);j++ {
				myNum = append(myNum, replace)
			}
		}
		var result string
		for  _,value := range myNum{

			result+=value
		}*/
		s = strings.Replace(s, word, replace, index-1)
	}
	return s, nil
}
