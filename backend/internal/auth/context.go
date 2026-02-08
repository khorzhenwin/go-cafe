package auth

import "context"

type contextKey string

const UserIDKey contextKey = "user_id"

func UserIDFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(UserIDKey).(uint)
	return id, ok
}
