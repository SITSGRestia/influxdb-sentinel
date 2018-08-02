package util

import (
	"time"
	"strconv"
	"per.zdh.org/influxdb-sentinel/global"
)


//格式字符串转时间戳
//time.ParseInLocation("20060102150405","20180223100000",time.Local)
func ConvertStringToTimeStamp(layout, value string, location *time.Location) (int64, error) {
	t, err := time.ParseInLocation(layout, value, location)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

//格式字符串转Time
//time.ParseInLocation("20060102150405","20180223100000",time.Local)
func ConvertStringToTime(layout, value string, location *time.Location) (*time.Time, error) {
	t, err := time.ParseInLocation(layout, value, location)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// 获取月初时间
func GetEarlyMonthUnix() int64 {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return tm.Unix()
}

// 获取零时时间
func GetZeroHourUnix() int64 {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return tm.Unix()
}

// 获取当前小时时间
func GetNowHourUnix() int64 {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	return tm.Unix()
}

// 获取当前时间
func GetNowUnix() int64 {
	return time.Now().Unix()
}

// 获取年初时间
func GetEarlyYearUnix() int64 {
	now := time.Now()
	tm := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	return tm.Unix()
}

func GetUnixToFormatString(timestamp int64, f string) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format(f)
}

func GetUnixToString(timestamp int64) string {
	return GetUnixToFormatString(timestamp, "2006-01-02 00:00:00")
}

func GetUnixToHourString(timestamp int64) string {
	return GetUnixToFormatString(timestamp, "15:04")
}

func GetUnixToDay(timestamp int64) int {
	tm := time.Unix(timestamp, 0)
	return tm.Day()
}

func GetDayString() string {
	now := time.Now().Unix()
	return GetUnixToFormatString(now, "20060102")
}
func GetUnixTimeString() string {
	now := time.Now().Unix()
	return GetUnixToFormatString(now, "20060102150304")
}

func GetUnixToDayTime(timestamp int64) string {
	month := GetUnixToMonth(timestamp)
	day := GetUnixToDay(timestamp)
	d := month + "." + strconv.Itoa(day)
	return d
}

func GetUnixToMonth(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return global.Month[tm.Month().String()]
}
