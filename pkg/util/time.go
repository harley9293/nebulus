package util

import "time"

func IsSameDay(time1, time2 int64) bool {
	return StartTimestamp(time1) == StartTimestamp(time2)
}

func TodayStartTimestamp() int64 {
	return StartTimestamp(time.Now().Unix())
}

func StartTimestamp(timestamp int64) int64 {
	timeStr := time.Unix(timestamp, 0).Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix()
}
