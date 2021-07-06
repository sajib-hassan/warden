package authorize

import (
	"context"
	"net/http"

	"github.com/go-chi/render"

	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// RequiredUser middleware restricts access to accounts having valid user in their jwt claims.
func RequiredUser(auth AuthStorer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			claims := jwt.ClaimsFromCtx(r.Context())
			user, err := auth.GetUser(claims.ID)
			if err != nil {
				// user deleted while access token still valid
				render.Render(w, r, ErrUnauthorized(err))
				return
			}
			ctx := context.WithValue(r.Context(), ctxUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
