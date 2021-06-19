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
	ctxAccount ctxKey = iota
	ctxProfile
)

// API provides application resources and handlers.
type API struct {
	Account *AccountResource
	Profile *ProfileResource
}

// NewAPI configures and returns application API.
func NewAPI() (*API, error) {
	accountStore := repos.NewUserStore()
	account := NewAccountResource(accountStore)

	profileStore := repos.NewProfileStore()
	profile := NewProfileResource(profileStore)

	api := &API{
		Account: account,
		Profile: profile,
	}
	return api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/account", a.Account.router())
	r.Mount("/profile", a.Profile.router())

	return r
}

func log() logrus.FieldLogger {
	//return logging.GetLogEntry(r)
	return logging.Logger
}
