package middleware

import (
	"context"
	"fmt"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func ContextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// RequestIDKey returns the context key for use in tests.
func RequestIDKey() string { return fmt.Sprintf("%v", requestIDKey) }