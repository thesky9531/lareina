package time

import (
	"context"
	"fmt"
	"time"
)

const (
	TimeIntervalMin uint32 = iota + 1
	TimeIntervalHour
	TimeIntervalDay
	TimeIntervalTwoDay
	TimeIntervalMonth
	TimeIntervalQuarter
	TimeIntervalYear
)

func GetTimeInterval(str string) uint32 {
	switch str {
	case "min":
		return TimeIntervalMin
	case "hour":
		return TimeIntervalHour
	case "day":
		return TimeIntervalDay
	case "two_day":
		return TimeIntervalTwoDay
	case "month":
		return TimeIntervalMonth
	case "quarter":
		return TimeIntervalQuarter
	case "year":
		return TimeIntervalYear
	default:
		return 0
	}
}

func UpdateContextDuration(d time.Duration, c context.Context) (time.Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if cTimeOut := time.Until(deadline); cTimeOut < d {
			// deliver small timeout
			return cTimeOut, c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, d)
	return d, ctx, cancel
}

// 获取某日开始（0分0秒等）
func GetOnedayStart(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取某日结束（23点59分等）
func GetOnedayEnd(d time.Time) time.Time {
	rel := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	rel = rel.AddDate(0, 0, 1)
	return rel.Add(-1)
}

// 获取某周一开始（0分0秒等）
func GetOneWeekStart(d time.Time) time.Time {
	offset := int(time.Monday - d.Weekday())
	if offset > 0 {
		offset = -6
	}
	d = d.AddDate(0, 0, offset)
	return GetOnedayStart(d)
}

// 获取某周末结束（23点59分等）
func GetOneWeekEnd(d time.Time) time.Time {
	offset := int(time.Sunday - d.Weekday())
	if offset == 0 {
		offset = -7
	}
	d = d.AddDate(0, 0, offset+7)
	return GetOnedayEnd(d)
}

// 获取某月1日开始（0分0秒等）
func GetOneMonthStart(d time.Time) time.Time {
	rel := GetOnedayStart(d)
	rel = rel.AddDate(0, 0, -rel.Day()+1)
	return rel
}

// 获取某月末结束（23点59分等）
func GetOneMonthEnd(d time.Time) time.Time {
	rel := GetOneMonthStart(d)
	rel = rel.AddDate(0, 1, 0)
	return rel.Add(-1)
}

// 获取某季度1日开始（0分0秒等）
func GetOneSeasonStart(d time.Time) time.Time {
	offset := int(d.Month()-time.January) % 3
	rel := GetOneMonthStart(d)
	rel = rel.AddDate(0, -offset, 0)
	return rel
}

// 获取某季度末结束（23点59分等）
func GetOneSeasonEnd(d time.Time) time.Time {
	rel := GetOneSeasonStart(d)
	rel = rel.AddDate(0, 4, 0)
	return rel.Add(-1)
}

// 获取某年开始（0分0秒等）
func GetOneYearStart(d time.Time) time.Time {
	rel := GetOneMonthStart(d)
	rel = rel.AddDate(0, int(time.January-rel.Month()), 0)
	return rel
}

// 获取某年结束（23点59分等）
func GetOneYearEnd(d time.Time) time.Time {
	rel := GetOneYearStart(d)
	rel = rel.AddDate(1, 0, 0)
	return rel.Add(-1)
}

// 获取某日开始（0分0秒等）
func GetOnedayStartString(d time.Time) string {
	return GetOnedayStart(d).Local().Format(time.RFC3339)
}

// 获取某日结束（23点59分等）
func GetOnedayEndString(d time.Time) string {
	return GetOnedayEnd(d).Local().Format(time.RFC3339)
}

// 获取某周一开始（0分0秒等）
func GetOneWeekStartString(d time.Time) string {
	return GetOneWeekStart(d).Local().Format(time.RFC3339)
}

// 获取某周末结束（23点59分等）
func GetOneWeekEndString(d time.Time) string {
	return GetOneWeekEnd(d).Local().Format(time.RFC3339)
}

// 获取某月1日开始（0分0秒等）
func GetOneMonthStartString(d time.Time) string {
	return GetOneMonthStart(d).Local().Format(time.RFC3339)
}

// 获取某月末结束（23点59分等）
func GetOneMonthEndString(d time.Time) string {
	return GetOneMonthEnd(d).Local().Format(time.RFC3339)
}

// 获取某季度1日开始（0分0秒等）
func GetOneSeasonStartString(d time.Time) string {
	return GetOneSeasonStart(d).Local().Format(time.RFC3339)
}

// 获取某季度末结束（23点59分等）
func GetOneSeasonEndString(d time.Time) string {
	return GetOneSeasonEnd(d).Local().Format(time.RFC3339)
}

// 获取某年开始（0分0秒等）
func GetOneYearStartString(d time.Time) string {
	return GetOneYearStart(d).Local().Format(time.RFC3339)
}

// 获取某年结束（23点59分等）
func GetOneYearEndString(d time.Time) string {
	return GetOneYearEnd(d).Local().Format(time.RFC3339)
}

func GetTimeLay(start *time.Time) (jj string) {

	timeout := fmt.Sprintf("%v", start)
	formatTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", timeout)
	jj = formatTime.Format("2006-01-02 15:04:05")
	return
}
