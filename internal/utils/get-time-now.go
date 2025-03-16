package utils

import "time"

func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetTimeNowWithoutSpacing() string {
	return time.Now().Format("20060102150405")
}
