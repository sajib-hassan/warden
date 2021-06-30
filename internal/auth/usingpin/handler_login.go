package usingpin

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/auth/mfa"
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
		twoFa, err := sentLoginOTP(rs.Store, acc, body)
		if err != nil {
			log().WithFields(logrus.Fields{
				"mobile":    body.Mobile,
				"device_id": body.DeviceId,
			}).Error(err)
			render.Render(w, r, ErrUnauthorized(err))
			return
		}

		mfa.SetChallengeHeader(w, twoFa)
		render.Respond(w, r, &loginOTPRequiredResponse{
			ChallengeRequired: true,
		})
		return
	}

	performLogin(w, r, acc, rs)
	return
}
