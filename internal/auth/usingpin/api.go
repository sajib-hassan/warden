package usingpin

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/logging"
)

// Resource implements PIN based user authentication against a database.
type Resource struct {
	TokenAuth *jwt.TokenAuth
	Store     AuthStorer
}

// NewResource returns a configured authentication resource.
func NewResource(authStore AuthStorer) (*Resource, error) {
	tokenAuth, err := jwt.NewTokenAuth()
	if err != nil {
		return nil, err
	}

	resource := &Resource{
		TokenAuth: tokenAuth,
		Store:     authStore,
	}

	resource.choresTicker()

	return resource, nil
}

// Router provides necessary routes for PIN based authentication flow.
func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/login", rs.login)
	r.Group(func(r chi.Router) {
		r.Use(rs.TokenAuth.Verifier())
		r.Use(jwt.AuthenticateRefreshJWT)
		r.Post("/refresh", rs.refresh)
		r.Post("/logout", rs.logout)
	})
	return r
}

func log() logrus.FieldLogger {
	return logging.Logger
}
