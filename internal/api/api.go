package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"github.com/sajib-hassan/warden/internal/app"
	"github.com/sajib-hassan/warden/internal/auth/usingpin"
	repos2 "github.com/sajib-hassan/warden/internal/db/repos"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/dbconn"
	"github.com/sajib-hassan/warden/pkg/logging"
)

// New configures application resources and routes.
func New() (*chi.Mux, error) {

	router, logger := InitAndBindRouter()

	err := dbconn.Connect()
	if err != nil {
		logger.WithField("module", "mongodb connect").Error(err)
		return nil, err
	}

	authStore := repos2.NewAuthStore()
	authResource, err := usingpin.NewResource(authStore)
	if err != nil {
		logger.WithField("module", "auth").Error(err)
		return nil, err
	}

	appAPI, err := app.NewAPI()
	if err != nil {
		logger.WithField("module", "app").Error(err)
		return nil, err
	}

	router.Mount("/auth", authResource.Router())
	router.Group(func(router chi.Router) {
		router.Use(authResource.TokenAuth.Verifier())
		router.Use(jwt.Authenticator)
		router.Mount("/api", appAPI.Router())
	})

	return router, nil
}

func InitAndBindRouter() (*chi.Mux, *logrus.Logger) {
	logger := logging.NewLogger()

	router := chi.NewRouter()
	router.Use(Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//router.Use(middleware.DefaultCompress)
	router.Use(middleware.Timeout(15 * time.Second))
	router.Use(middleware.Heartbeat("/ping"))

	router.Use(logging.NewStructuredLogger(logger))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := newAPIError("Not Found", errURINotFound, nil)
		res := response{
			Code:   http.StatusNotFound,
			Errors: apiErrors{*err},
		}
		res.serveJSON(w)
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		err := newAPIError("Method not allowed", errInvalidMethod, nil)
		res := response{
			Code:   http.StatusMethodNotAllowed,
			Errors: apiErrors{*err},
		}
		res.serveJSON(w)
	})

	return router, logger
}
