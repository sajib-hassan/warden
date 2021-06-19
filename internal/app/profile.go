package app

import (
	"context"
	"errors"
	"net/http"

	models2 "github.com/sajib-hassan/warden/internal/db/models"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// The list of error types returned from account resource.
var (
	ErrProfileValidation = errors.New("profile validation error")
)

// ProfileStore defines database operations for a profile.
type ProfileStore interface {
	Get(accountID int) (*models2.Profile, error)
	Update(p *models2.Profile) error
}

// ProfileResource implements profile management handler.
type ProfileResource struct {
	Store ProfileStore
}

// NewProfileResource creates and returns a profile resource.
func NewProfileResource(store ProfileStore) *ProfileResource {
	return &ProfileResource{
		Store: store,
	}
}

func (rs *ProfileResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(rs.profileCtx)
	r.Get("/", rs.get)
	r.Put("/", rs.update)
	return r
}

func (rs *ProfileResource) profileCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := jwt.ClaimsFromCtx(r.Context())
		p, err := rs.Store.Get(claims.ID)
		if err != nil {
			log().WithField("profileCtx", claims.Sub).Error(err)
			render.Render(w, r, ErrInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), ctxProfile, p)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type profileRequest struct {
	*models2.Profile
	ProtectedID int `json:"id"`
}

func (d *profileRequest) Bind(r *http.Request) error {
	return nil
}

type profileResponse struct {
	*models2.Profile
}

func newProfileResponse(p *models2.Profile) *profileResponse {
	return &profileResponse{
		Profile: p,
	}
}

func (rs *ProfileResource) get(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(ctxProfile).(*models2.Profile)
	render.Respond(w, r, newProfileResponse(p))
}

func (rs *ProfileResource) update(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(ctxProfile).(*models2.Profile)
	data := &profileRequest{Profile: p}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}

	if err := rs.Store.Update(p); err != nil {
		switch err.(type) {
		case validation.Errors:
			render.Render(w, r, ErrValidation(ErrProfileValidation, err.(validation.Errors)))
			return
		}
		render.Render(w, r, ErrRender(err))
		return
	}
	render.Respond(w, r, newProfileResponse(p))
}
