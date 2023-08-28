package uf

import "time"

var OneDayMilli = 86400000
var OneDaySecs = 86400

var EmptyDate = 9999999999999

func TimeFromYm(year int, month int) time.Time {
	return time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
}

func StampFromYm(year int, month int) int64 {
	return TimeFromYm(year, month).UnixMilli()
}

func Now() time.Time {
	return time.Now().UTC()
}

func NowStamp() int64 {
	return time.Now().UnixMilli()
}

func InXDays(days int) time.Time {
	return Now().AddDate(0, 0, days)
}

func InXMinutesStamp(minutes int) int64 {
	return NowStamp() + int64(minutes*60000)
}

func InXDaysStamp(days int) int64 {
	return InXDays(days).UnixMilli()
}
