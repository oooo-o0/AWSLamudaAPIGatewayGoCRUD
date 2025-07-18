package auth

import (
	"context"
)

type contextKey string

const userClaimsKey = contextKey("userClaims")

func SetUserClaimsCtx(ctx context.Context, claims *CustomClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

func GetUserClaimsCtx(ctx context.Context) (*CustomClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*CustomClaims)
	return claims, ok
}
