package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sajib-hassan/warden/internal/app"
	"github.com/sajib-hassan/warden/internal/repos"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/auth/pwdless"
	"github.com/sajib-hassan/warden/pkg/dbconn"
	"github.com/sajib-hassan/warden/pkg/email"
	"github.com/sajib-hassan/warden/pkg/logging"
	"net/http"
	"time"
)

// New configures application resources and routes.
func New() (*chi.Mux, error) {
	logger := logging.NewLogger()

	db, err := dbconn.Connect()
	if err != nil {
		logger.WithField("module", "database").Error(err)
		return nil, err
	}

	mailer, err := email.NewMailer()
	if err != nil {
		logger.WithField("module", "email").Error(err)
		return nil, err
	}

	authStore := repos.NewAuthStore(db)
	authResource, err := pwdless.NewResource(authStore, mailer)
	if err != nil {
		logger.WithField("module", "auth").Error(err)
		return nil, err
	}

	//adminAPI, err := admin.NewAPI(db)
	//if err != nil {
	//	logger.WithField("module", "admin").Error(err)
	//	return nil, err
	//}

	appAPI, err := app.NewAPI(db)
	if err != nil {
		logger.WithField("module", "app").Error(err)
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP)
	//router.Use(middleware.DefaultCompress)
	router.Use(middleware.Timeout(15 * time.Second))

	router.Use(logging.NewStructuredLogger(logger))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// @todo Sajib
		//err := newApiError("Not Found", errURINotFound, nil)
		//panic(err)
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		// @todo Sajib
		//err := newApiError("Method not allowed", errInvalidMethod, nil)
		//res := response{
		//	Code:   http.StatusMethodNotAllowed,
		//	Errors: apiErrors{*err},
		//}
		//res.serveJSON(w)
	})

	router.Mount("/auth", authResource.Router())
	router.Group(func(router chi.Router) {
		router.Use(authResource.TokenAuth.Verifier())
		router.Use(jwt.Authenticator)
		//router.Mount("/admin", adminAPI.Router())
		router.Mount("/api", appAPI.Router())
	})
	//
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return router, nil
}
