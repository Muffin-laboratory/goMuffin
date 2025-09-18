package routes

import (
	"net/http"
	"time"

	"git.wh64.net/muffin/goMuffin/configs"
	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Data: map[string]any{
			"happy":   "hacking!",
			"date":    time.Now(),
			"version": configs.MuffinVersion,
			"branch":  configs.CurrentBranch,
		}},
	)
}
