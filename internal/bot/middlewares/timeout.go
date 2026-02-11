package middlewares

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/handler"
)

func TimeoutMiddleware(timeout time.Duration) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(inter *handler.InteractionEvent) error {
			ctx, cancel := context.WithTimeout(inter.Ctx, timeout)
			defer cancel()

			inter.Ctx = ctx

			return next(inter)
		}
	}
}
