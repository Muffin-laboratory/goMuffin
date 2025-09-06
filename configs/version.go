package configs

import (
	"strconv"
	"time"

	"git.wh64.net/muffin/goMuffin/utils"
)

const MuffinVersion = "7.0.0-madeleine_v_two_develop.250906a"

var updatedString string = utils.RegexpDecimals.FindAllStringSubmatch(MuffinVersion, -1)[3][0]

var UpdatedAt *time.Time = func() *time.Time {
	year, _ := strconv.Atoi("20" + updatedString[0:2])
	monthInt, _ := strconv.Atoi(updatedString[2:4])
	month := time.Month(monthInt)
	day, _ := strconv.Atoi(updatedString[4:6])
	time := time.Date(year, month, day, 0, 0, 0, 0, &time.Location{})
	return &time
}()
