package util

import (
	"testing"
	"time"
)

func TestIsSameDay(t *testing.T) {
	if !IsSameDay(time.Now().Unix(), time.Now().Add(1*time.Hour).Unix()) {
		t.Error("IsSameDay(Now(), Now())")
	}

	if IsSameDay(time.Now().Unix(), time.Now().Add(24*time.Hour).Unix()) {
		t.Error("IsSameDay(Now(), Now()+24h)")
	}
}

func TestTodayStartTimestamp(t *testing.T) {
	if TodayStartTimestamp() != StartTimestamp(time.Now().Unix()) {
		t.Error("TodayStartTimestamp() != StartTimestamp(Now())")
	}
}
