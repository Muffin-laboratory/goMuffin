package routes

import (
	"net/http"

	"git.wh64.net/muffin/goMuffin/chatbot"
	"github.com/labstack/echo/v4"
)

func ReloadPrompt(c echo.Context) error {
	if err := chatbot.GetChatBot().ReloadPrompt(); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Message: "Reloading prompt is success.",
	})
}

func GetPrompt(c echo.Context) error {
	prompt := chatbot.GetChatBot().GetPrompt()

	return c.JSON(http.StatusOK, Response{
		Data: prompt,
	})
}
