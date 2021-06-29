package usingpin

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

func (rs *Resource) logout(w http.ResponseWriter, r *http.Request) {
	rt := jwt.RefreshTokenFromCtx(r.Context())
	token, err := rs.Store.GetToken(rt)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(jwt.ErrTokenExpired))
		return
	}
	rs.Store.DeleteToken(token)

	render.Respond(w, r, http.NoBody)
}
