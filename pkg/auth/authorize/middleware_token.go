package authorize

import (
	"context"
	"net/http"

	"github.com/go-chi/render"

	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// RequiredToken middleware restricts access to accounts having valid token in their jwt claims.
func RequiredToken(auth AuthStorer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			claims := jwt.ClaimsFromCtx(r.Context())
			token, err := auth.GetToken(claims.Token)
			if err != nil {
				// token deleted while access token still valid
				render.Render(w, r, ErrUnauthorized(jwt.ErrInvalidAccessToken))
				return
			}
			ctx := context.WithValue(r.Context(), ctxToken, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
