package contextutils

import "context"

type xUserId string

const xUserIdKey xUserId = "x-user-id"

func ContextWithXUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, xUserIdKey, userId)
}

func XUserIdFromContext(ctx context.Context) string {
	userIdAny := ctx.Value(xUserIdKey)
	if userId, ok := userIdAny.(string); ok {
		return userId
	}
	return ""
}
