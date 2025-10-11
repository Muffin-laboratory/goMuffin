package configs

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"git.wh64.net/muffin/goMuffin/utils"
)

var MuffinVersion = fmt.Sprintf("7.0.0-pretzel_%s.251011b", CurrentBranch)

var updatedString string = utils.RegexpDecimals.FindAllStringSubmatch(MuffinVersion, -1)[3][0]

var UpdatedAt *time.Time = func() *time.Time {
	year, _ := strconv.Atoi("20" + updatedString[0:2])
	monthInt, _ := strconv.Atoi(updatedString[2:4])
	month := time.Month(monthInt)
	day, _ := strconv.Atoi(updatedString[4:6])
	time := time.Date(year, month, day, 0, 0, 0, 0, &time.Location{})
	return &time
}()

var CurrentBranch = func() string {
	var out strings.Builder

	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "release"
	}

	if strings.Contains(out.String(), "main") {
		return "release"
	}

	return strings.Trim(out.String(), "\n")
}()
