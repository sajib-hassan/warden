package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/sajib-hassan/warden/internal/auth/usingpin"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// The list of error types returned from user resource.
var (
	ErrUserValidation = errors.New("user validation error")
)

// UserResource implements user management handler.
type UserResource struct {
	Store UserStore
}

// NewUserResource creates and returns an user resource.
func NewUserResource(store UserStore) *UserResource {
	return &UserResource{
		Store: store,
	}
}

func (rs *UserResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(rs.userCtx)
	r.Get("/", rs.get)
	r.Put("/", rs.update)
	r.Delete("/", rs.delete)
	r.Route("/token/{tokenID}", func(r chi.Router) {
		r.Put("/", rs.updateToken)
		r.Delete("/", rs.deleteToken)
	})
	return r
}

func (rs *UserResource) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := jwt.ClaimsFromCtx(r.Context())
		log().WithField("user_id", claims.ID)
		user, err := rs.Store.Get(claims.ID)
		if err != nil {
			// user deleted while access token still valid
			render.Render(w, r, ErrUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type userRequest struct {
	*usingpin.User
	// override protected data here, although not really necessary here
	// as we limit updated database columns in store as well
	ProtectedID     int      `json:"id"`
	ProtectedActive bool     `json:"active"`
	ProtectedRoles  []string `json:"roles"`
}

func (d *userRequest) Bind(r *http.Request) error {
	// d.ProtectedActive = true
	// d.ProtectedRoles = []string{}
	return nil
}

type userResponse struct {
	*usingpin.User
}

func newUserResponse(a *usingpin.User) *userResponse {
	resp := &userResponse{User: a}
	return resp
}

func (rs *UserResource) get(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxUser).(*usingpin.User)
	render.Respond(w, r, newUserResponse(acc))
}

func (rs *UserResource) update(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxUser).(*usingpin.User)
	data := &userRequest{User: acc}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.Store.Update(acc); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrUserValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newUserResponse(acc))
}

func (rs *UserResource) delete(w http.ResponseWriter, r *http.Request) {
	acc := r.Context().Value(ctxUser).(*usingpin.User)
	if err := rs.Store.Delete(acc); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, http.NoBody)
}

type tokenRequest struct {
	Identifier  string
	ProtectedID int `json:"id"`
}

func (d *tokenRequest) Bind(r *http.Request) error {
	d.Identifier = strings.TrimSpace(d.Identifier)
	return nil
}

func (rs *UserResource) updateToken(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "tokenID")
	data := &tokenRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	acc := r.Context().Value(ctxUser).(*usingpin.User)
	for _, t := range acc.Token {
		if t.ID.Hex() == id {
			jt := &jwt.Token{
				Identifier: data.Identifier,
			}
			jt.SetID(t.ID)
			if err := rs.Store.UpdateToken(jt); err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
		}
	}
	render.Respond(w, r, http.NoBody)
}

func (rs *UserResource) deleteToken(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "tokenID")
	acc := r.Context().Value(ctxUser).(*usingpin.User)
	for _, t := range acc.Token {
		if t.ID.Hex() == id {
			jt := &jwt.Token{}
			jt.SetID(t.ID)
			rs.Store.DeleteToken(jt)
		}
	}
	render.Respond(w, r, http.NoBody)
}
