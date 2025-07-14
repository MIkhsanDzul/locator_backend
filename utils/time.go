package utils

import (
	"time"
	// "fmt"
)

// Dari UnixMilli ke string waktu yang mudah dibaca
func FormatTimestamp(millis int64) string {
	t := time.UnixMilli(millis)
	return t.Format("2006-01-02 15:04:05")
}

// (Opsional) Format dengan zona waktu lokal
func FormatTimestampWithZone(millis int64, loc string) string {
	location, err := time.LoadLocation(loc)
	if err != nil {
		location = time.UTC
	}
	t := time.UnixMilli(millis).In(location)
	return t.Format("2006-01-02 15:04:05 MST")
}