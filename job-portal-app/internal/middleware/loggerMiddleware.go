package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type key string

const TraceIDKey key = "1"

func (m *Mid) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidStr := uuid.NewString()

		ctx := c.Request.Context()

		ctx = context.WithValue(ctx, TraceIDKey, uuidStr)

		c.Request = c.Request.WithContext(ctx)

		log.Info().Str("Trace ID", uuidStr).Str("Method", c.Request.Method).Str("URL Path", c.Request.URL.Path).Msg("Request Started")

		defer log.Info().Str("Trace ID", uuidStr).Str("Method", c.Request.Method).Str("URL Path", c.Request.URL.Path).Int("Status Code", c.Writer.Status()).Msg("Request processing completed")

		c.Next()
	}
}
