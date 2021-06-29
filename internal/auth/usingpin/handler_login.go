package usingpin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/validator"
)

type loginRequest struct {
	Mobile   string `json:"mobile"`
	Pin      string `json:"pin"`
	DeviceId string `json:"device_id"`
}

type loginResponse struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
	Slug    string `json:"user_slug"`
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
}

// "user_slug": userObj.ID,
//			"name":      userObj.FirstName + " " + userObj.LastName,
//			"phone":     userObj.Phone,

type loginOTPRequiredResponse struct {
	ChallengeRequired bool `json:"challenge_required"`
}

func (body *loginRequest) Bind(_ *http.Request) error {
	body.Mobile = strings.TrimSpace(body.Mobile)
	body.Mobile = strings.ToLower(body.Mobile)

	body.Pin = strings.TrimSpace(body.Pin)
	body.Pin = strings.ToLower(body.Pin)

	body.DeviceId = strings.TrimSpace(body.DeviceId)

	loginPinLength := viper.GetInt("auth_login_pin_length")

	return validation.ValidateStruct(body,
		validation.Field(&body.Mobile,
			validation.Required,
			validation.By(validator.CheckValidBDMobileNumber),
		),
		validation.Field(&body.Pin, validation.Required, is.Digit, validation.Length(loginPinLength, loginPinLength)),
		validation.Field(&body.DeviceId, validation.Required, is.ASCII),
	)
}

func (rs *Resource) login(w http.ResponseWriter, r *http.Request) {
	body := &loginRequest{}
	if err := render.Bind(r, body); err != nil {
		log().WithFields(logrus.Fields{
			"mobile":    body.Mobile,
			"pin":       "*******",
			"device_id": body.DeviceId,
		}).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidLogin))
		return
	}

	acc, err := rs.Store.GetUserByMobile(body.Mobile)
	if err != nil {
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
			"mobile":    body.Mobile,
			"pin":       "*******",
			"device_id": body.DeviceId,
		}).Warn(err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidLogin))
		return
	}

	device, err := rs.Store.GetTrustedDevice(acc.ID.Hex(), body.DeviceId)
	if err != nil {
		log().WithField("device_id", body.DeviceId).Error(err)
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	if device == nil {
		twoFaToken, err := SentLoginOTP(rs.Store, acc)
		if err != nil {
			log().WithFields(logrus.Fields{
				"mobile":    body.Mobile,
				"device_id": body.DeviceId,
			}).Error(err)
			render.Render(w, r, ErrUnauthorized(err))
			return
		}

		w.Header().Set("Pathao-Challenge-Token", twoFaToken)
		render.Respond(w, r, &loginOTPRequiredResponse{
			ChallengeRequired: true,
		})
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
		Slug:    acc.ID.Hex(),
		Name:    acc.Name,
		Mobile:  acc.Mobile,
	})
}
