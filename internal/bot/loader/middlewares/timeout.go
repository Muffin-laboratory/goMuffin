package middlewares

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
)

func Timeout(timeout time.Duration) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(inter *handler.InteractionEvent) error {
			ctx, cancel := context.WithTimeout(inter.Ctx, timeout)
			defer cancel()

			inter.Ctx = ctx

			return next(inter)
		}
	}
}

func TimeoutAndDefer(timeout time.Duration, iType discord.InteractionType, updateMessage, ephemeral bool) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		next = middleware.Defer(iType, updateMessage, ephemeral)(next)
		return Timeout(timeout)(next)
	}
}
