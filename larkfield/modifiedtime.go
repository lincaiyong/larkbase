package larkfield

import (
	"time"
)

const ModifiedTimeFieldName = "modified_time"

type ModifiedTime int

func TimeToModifiedTime(t time.Time) ModifiedTime {
	t = t.In(beijingTZ)
	yy := t.Year() % 2000
	mm := int(t.Month())
	dd := t.Day()
	hh := t.Hour()
	mn := t.Minute()
	return ModifiedTime(yy*1e8 + mm*1e6 + dd*1e4 + hh*1e2 + mn)
}

func ModifiedTimeToTime(m ModifiedTime) time.Time {
	i := int(m)
	yy := i/1e8 + 2000
	i = i % 1e8
	mm := i / 1e6
	i = i % 1e6
	dd := i / 1e4
	i = i % 1e4
	hh := i / 1e2
	mn := i % 1e2
	t := time.Date(yy, time.Month(mm), dd, hh, mn, 0, 0, beijingTZ)
	return t
}
