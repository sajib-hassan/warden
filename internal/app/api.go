// Package app ties together application resources and handlers.
package app

import (
	"github.com/sajib-hassan/warden/internal/db/repos"
	"github.com/sajib-hassan/warden/pkg/logging"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type ctxKey int

const (
	ctxUser ctxKey = iota
	ctxProfile
)

// API provides application resources and handlers.
type API struct {
	User    *UserResource
	Profile *ProfileResource
}

// NewAPI configures and returns application API.
func NewAPI() (*API, error) {
	userStore := repos.NewUserStore()
	user := NewUserResource(userStore)

	profileStore := repos.NewProfileStore()
	profile := NewProfileResource(profileStore)

	api := &API{
		User:    user,
		Profile: profile,
	}
	return api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/v1/user", a.User.router())
	r.Mount("/v1/profile", a.Profile.router())

	return r
}

func log() logrus.FieldLogger {
	//return logging.GetLogEntry(r)
	return logging.Logger
}
