package authorize

import (
	"context"

	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

type ctxKey int

const (
	ctxUser ctxKey = iota
	ctxToken
)
const (
	TwoFaLoginServiceType = "*authorize.User"
	TwoFaLoginFor         = "login"
)

// CurrentUserFromCtx retrieves the current user from request context.
func CurrentUserFromCtx(ctx context.Context) *User {
	return ctx.Value(ctxUser).(*User)
}

// CurrentTokenFromCtx retrieves the token from context.
func CurrentTokenFromCtx(ctx context.Context) *jwt.Token {
	return ctx.Value(ctxToken).(*jwt.Token)
}
