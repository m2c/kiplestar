package httptool

import "bytes"

// letter is upper or not
func IsUpperLetter(letter rune) bool {
	if letter >= 'A' && letter <= 'Z' {
		return true
	} else {
		return false
	}
}

func IsLowerLetter(letter rune) bool {
	if letter >= 'a' && letter <= 'z' {
		return true
	} else {
		return false
	}
}

func TransLetterToUpper(letter rune) string {
	if IsLowerLetter(letter) {
		letter -= 'a' - 'A'
	}
	return string(letter)
}

func TransLetterToLower(letter rune) string {
	if IsUpperLetter(letter) {
		letter += 'a' - 'A'
	}
	return string(letter)
}

// like transform "to_lower_snake_case" to "toLowerSnakeCase"
func ToLowerCamelCase(s string) string {
	var dst bytes.Buffer
	var flag bool
	for index, letter := range s {
		if index == 0 {
			dst.WriteString(TransLetterToLower(letter))
		} else if letter == '_' || letter == '-' {
			flag = true
		} else if flag {
			flag = false
			dst.WriteString(TransLetterToUpper(letter))
		} else {
			dst.WriteString(string(letter))
		}
	}

	return dst.String()
}
