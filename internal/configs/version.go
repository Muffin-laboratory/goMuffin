package configs

import (
	"strconv"
	"time"
)

var (
	MuffinVersion = "0.0.0-local_debug.0002110000.000211"
	updatedString = "0002110000"
)

func UpdatedAt() *time.Time {
	year, _ := strconv.Atoi("20" + updatedString[0:2])
	monthInt, _ := strconv.Atoi(updatedString[2:4])
	month := time.Month(monthInt)
	day, _ := strconv.Atoi(updatedString[4:6])
	hour, _ := strconv.Atoi(updatedString[6:8])
	minute, _ := strconv.Atoi(updatedString[8:10])
	time := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	return &time
}
