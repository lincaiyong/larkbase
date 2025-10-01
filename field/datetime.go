package field

import (
	"time"
)

const dateTimeLayout = "2006-01-02 15:04:05"

var beijingTZ = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("Beijing", 8*60*60)
	}
	return loc
}()

func beijingDateTimeStrToUnixSeconds(s string) (int64, error) {
	t, err := time.ParseInLocation(dateTimeLayout, s, beijingTZ)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func beijingDateTimeStrToTime(s string) (time.Time, error) {
	return time.ParseInLocation(dateTimeLayout, s, beijingTZ)
}

func unixSecondsToBeijingDateTimeStr(timestamp int64) string {
	return time.Unix(timestamp, 0).In(beijingTZ).Format(dateTimeLayout)
}

func utcDateTimeStrToUnixSeconds(s string) (int64, error) {
	t, err := time.ParseInLocation(dateTimeLayout, s, time.UTC)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func unixSecondsToUTCDateTimeStr(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(dateTimeLayout)
}
