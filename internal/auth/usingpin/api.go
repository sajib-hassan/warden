package usingpin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/logging"
)

// AuthStorer defines database operations on users and tokens.
type AuthStorer interface {
	GetUser(id string) (*User, error)
	GetUserByMobile(mobile string) (*User, error)
	UpdateUser(a *User) error

	GetToken(token string) (*jwt.Token, error)
	CreateOrUpdateToken(t *jwt.Token) error
	DeleteToken(t *jwt.Token) error
	PurgeExpiredToken() error
}

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

type loginRequest struct {
	Mobile string
	Pin    string
}

type loginResponse struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func (body *loginRequest) Bind(r *http.Request) error {
	body.Mobile = strings.TrimSpace(body.Mobile)
	body.Mobile = strings.ToLower(body.Mobile)

	body.Pin = strings.TrimSpace(body.Pin)
	body.Pin = strings.ToLower(body.Pin)

	loginPinLength := viper.GetInt("auth_login_pin_length")

	return validation.ValidateStruct(body,
		validation.Field(&body.Mobile,
			validation.Required,
			validation.By(CheckValidBDMobileNumber),
		),
		validation.Field(&body.Pin, validation.Required, is.Digit, validation.Length(loginPinLength, loginPinLength)),
	)
}

func (rs *Resource) login(w http.ResponseWriter, r *http.Request) {
	body := &loginRequest{}
	if err := render.Bind(r, body); err != nil {
		log().WithFields(logrus.Fields{
			"mobile": body.Mobile,
			"pin":    "*******",
		}).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidLogin))
		return
	}

	acc, err := rs.Store.GetUserByMobile(body.Mobile)
	if err != nil || acc == nil {
		log().WithField("mobile", body.Mobile).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrUnknownLogin))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, ErrUnauthorized(ErrLoginDisabled))
		return
	}

	if ok, err := acc.isPinMatched(body.Pin); !ok {
		log().WithFields(logrus.Fields{
			"mobile": body.Mobile,
			"pin":    "*******",
		}).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidLogin))
		return
	}

	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()

	token := &jwt.Token{
		Token:      uuid.Must(uuid.NewV4()).String(),
		Expiry:     time.Now().Add(rs.TokenAuth.JwtRefreshExpiry),
		UserID:     acc.ID.Hex(),
		Mobile:     ua.Mobile(),
		Identifier: fmt.Sprintf("%s on %s", browser, ua.OS()),
	}

	if err := rs.Store.CreateOrUpdateToken(token); err != nil {
		log().Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	access, refresh, err := rs.TokenAuth.GenTokenPair(acc.Claims(), token.Claims())
	if err != nil {
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
