package usingpin

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/gofrs/uuid"

	"github.com/sajib-hassan/warden/pkg/auth/authorize"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

type refreshResponse struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func (rs *Resource) refresh(w http.ResponseWriter, r *http.Request) {
	rt := jwt.RefreshTokenFromCtx(r.Context())

	token, err := rs.Store.GetToken(rt)
	if err != nil {
		render.Render(w, r, authorize.ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	if time.Now().After(token.Expiry) {
		rs.Store.DeleteToken(token)
		render.Render(w, r, authorize.ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	acc, err := rs.Store.GetUser(token.UserID)
	if err != nil {
		render.Render(w, r, authorize.ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, authorize.ErrUnauthorized(ErrLoginDisabled))
		return
	}

	srv := &service{w: w, r: r, rs: rs}

	token.Token = uuid.Must(uuid.NewV4()).String()
	token.Expiry = time.Now().Add(rs.TokenAuth.JwtRefreshExpiry)

	access, refresh, done := srv.getTokens(acc, token)
	if !done {
		return
	}

	render.Respond(w, r, &refreshResponse{
		Access:  access,
		Refresh: refresh,
	})
}
