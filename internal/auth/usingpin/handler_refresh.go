package usingpin

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/gofrs/uuid"

	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

func (rs *Resource) refresh(w http.ResponseWriter, r *http.Request) {
	rt := jwt.RefreshTokenFromCtx(r.Context())

	token, err := rs.Store.GetToken(rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	if time.Now().After(token.Expiry) {
		rs.Store.DeleteToken(token)
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}

	acc, err := rs.Store.GetUser(token.UserID)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	token.Token = uuid.Must(uuid.NewV4()).String()
	token.Expiry = time.Now().Add(rs.TokenAuth.JwtRefreshExpiry)
	token.UpdatedAt = time.Now()

	access, refresh, err := rs.TokenAuth.GenTokenPair(acc.Claims(), token.Claims())
	if err != nil {
		log().Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	if err := rs.Store.CreateOrUpdateToken(token); err != nil {
		log().Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	acc.LastLogin = time.Now()
	if err := rs.Store.UpdateUser(acc); err != nil {
		log().Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, &loginResponse{
		Access:  access,
		Refresh: refresh,
	})
}
