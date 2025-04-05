package common

import (
	"fmt"
	"strconv"
	"time"
)

const (
	Q1 = "Q1"
	Q2 = "Q2"
	Q3 = "Q3"
	Q4 = "Q4"
)

func ValidateTime(timeString string, timeZone *time.Location) (int, string, error) {
	var year int
	var quarter string
	// validate timestamp
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return year, quarter, fmt.Errorf("invalid timestamp")
	}
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, timeZone).Unix()
	end := time.Date(2099, 12, 31, 23, 59, 59, 0, timeZone).Unix()
	if timestamp < start || end < timestamp {
		return year, quarter, fmt.Errorf("the timestamp is invalid, not in range 2000 - 2099")
	}
	// get year and quarter information
	t := time.Unix(timestamp, 0).In(timeZone)
	year = t.Year()
	quarter = getQuarter(t.Month())

	return year, quarter, nil
}

func getQuarter(month time.Month) string {
	switch month {
	case time.January, time.February, time.March:
		return Q1
	case time.April, time.May, time.June:
		return Q2
	case time.July, time.August, time.September:
		return Q3
	default:
		return Q4
	}
}
