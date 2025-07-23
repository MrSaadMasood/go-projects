package contexter

import (
	"context"
)

type (
	ContextCounter int
	TokenUser      struct {
		Email string `json:"email"`
	}
)

const (
	userKey ContextCounter = iota
)

func ContextWithUser(ctx context.Context, user TokenUser) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (TokenUser, bool) {
	user, ok := ctx.Value(userKey).(TokenUser)
	return user, ok
}
