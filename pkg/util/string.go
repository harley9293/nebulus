package util

import (
	"strconv"
	"strings"
)

func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArray := []rune(str)
	if strArray[0] >= 97 && strArray[0] <= 122 {
		strArray[0] -= 32
	}
	return string(strArray)
}

func Unicode2Hans(raw []byte) []byte {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return raw
	}
	return []byte(str)
}
