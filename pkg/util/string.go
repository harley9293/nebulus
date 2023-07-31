package util

import (
	"strconv"
	"strings"
	"time"
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

func IsSameDay(time1, time2 int64, offset int) bool {
	return StartTimestamp(time1, offset) == StartTimestamp(time2, offset)
}

func TodayStartTimestamp(offset int) int64 {
	return StartTimestamp(time.Now().Unix(), offset)
}

func StartTimestamp(timestamp int64, offset int) int64 {
	timeStr := time.Unix(timestamp, 0).Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	t2 := t.Unix() + int64(offset)
	if t2 > timestamp {
		t2 -= 86400
	}
	return t2
}

func Unicode2Hans(raw []byte) []byte {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return raw
	}
	return []byte(str)
}
