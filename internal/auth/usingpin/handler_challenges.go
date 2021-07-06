package usingpin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"github.com/mssola/user_agent"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/auth/authorize"
	"github.com/sajib-hassan/warden/pkg/auth/mfa"
)

type validateOTPRequest struct {
	ChallengeCode   string `json:"challenge_code"`
	TrustThisDevice bool   `json:"trust_this_device"`
	Token           string `json:"token"`
}

func (body *validateOTPRequest) Bind(_ *http.Request) error {
	body.ChallengeCode = strings.TrimSpace(body.ChallengeCode)
	return validation.ValidateStruct(body,
		validation.Field(&body.TrustThisDevice, validation.Required),
		validation.Field(&body.Token, validation.Required, is.ASCII),
		validation.Field(&body.ChallengeCode, validation.Required, is.Digit, validation.Length(6, 6)),
	)
}

func (rs *Resource) validateOTP(w http.ResponseWriter, r *http.Request) {
	body := &validateOTPRequest{}
	body.Token = mfa.GetChallengeHeader(r)
	if err := render.Bind(r, body); err != nil {
		log().WithField("challenge_code", "******").Warn(err)
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrInvalidChallengeToken))
		return
	}

	twoFa, err := rs.Store.GetTwoFa(body.Token, authorize.TwoFaLoginServiceType, authorize.TwoFaLoginFor)
	if err != nil {
		log().Error(err)
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrInvalidChallengeToken))
		return
	}

	if twoFa.IsExpired() {
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrResendTimeExpired))
		return
	}

	if twoFa.IsExceedMaxRetry() {
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrMaxResendExceed))
		return
	}

	acc, err := rs.Store.GetUser(twoFa.ServiceId)
	if err != nil {
		log().Error(err)
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrChallengeTokenUnauthorized))
		return
	}

	if !acc.CanLogin() {
		render.Render(w, r, authorize.ErrUnauthorized(ErrLoginDisabled))
		return
	}

	srv := &service{w: w, r: r, rs: rs}

	vc := mfa.NewVerificationCode()

	if !vc.IsValidChallengeCode(body.ChallengeCode, acc.Secret, twoFa.Challenge) {
		twoFa.RetryCount = twoFa.RetryCount + 1
		if err := rs.Store.CreateOrUpdateTwoFa(twoFa); err != nil {
			log().Error(err)
		}
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrChallengeCodeMissMatched))
		return
	}

	if body.TrustThisDevice == true {
		var deviceId string
		var ok bool
		if x, found := twoFa.MetaData["device_id"]; found {
			if deviceId, ok = x.(string); !ok && len(deviceId) <= 0 {
				log().Error(errors.New("DeviceId not found in two-fa metadata"))
				render.Render(w, r, ErrInternalServerError)
			}
		}

		device, err := rs.Store.GetTrustedDevice(acc.ID.Hex(), deviceId)
		if err != nil {
			log().WithField("device_id", deviceId).Error(err)
			render.Render(w, r, ErrInternalServerError)
			return
		}

		if device != nil {
			device.IsAuthorized = true
		} else {
			ua := user_agent.New(r.UserAgent())
			name, version := ua.Engine()
			device = &authorize.Device{
				UserID:       acc.ID.Hex(),
				Identifier:   deviceId,
				Name:         fmt.Sprintf("%v %v", ua.OS(), ua.Platform()),
				IsAuthorized: true,
				Details: map[string]interface{}{
					"engine_name":    name,
					"engine_version": version,
				},
			}
		}
		rs.Store.RegisterAsTrustedDevice(device)
	}

	err = rs.Store.DeleteTwoFa(twoFa)
	if err != nil {
		log().Error(err)
		render.Render(w, r, ErrInternalServerError)
	}

	srv.performLogin(acc)
	return
}

type resendOTPResponse struct {
	Resend bool `json:"resend"`
}

func (rs *Resource) resendOTP(w http.ResponseWriter, r *http.Request) {
	token := mfa.GetChallengeHeader(r)
	if token == "" {
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrInvalidChallengeToken))
		return
	}

	twoFa, err := rs.Store.GetTwoFa(token, authorize.TwoFaLoginServiceType, authorize.TwoFaLoginFor)
	if err != nil {
		log().Error(err)
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrInvalidChallengeToken))
		return
	}

	if twoFa.IsExceedMaxResend() {
		rs.Store.DeleteTwoFa(twoFa)
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrMaxResendExceed))
		return
	}

	if twoFa.IsResendExpired() {
		rs.Store.DeleteTwoFa(twoFa)
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrResendTimeExpired))
		return
	}

	acc, err := rs.Store.GetUser(twoFa.ServiceId)
	if err != nil {
		log().Error(err)
		render.Render(w, r, authorize.ErrUnauthorized(mfa.ErrChallengeTokenUnauthorized))
		return
	}

	vc := mfa.NewVerificationCode()

	validUntil := time.Now().UTC().Add(viper.GetDuration("AUTH_LOGIN_OTP_TIMEOUT")).Unix()
	ResendValidUntil := time.Now().UTC().Add(viper.GetDuration("AUTH_LOGIN_OTP_RESEND_TIMEOUT")).Unix()
	twoFa.RetryCount = 0
	twoFa.ResendCount = twoFa.ResendCount + 1
	twoFa.Challenge = vc.Hashed(acc.Secret)
	twoFa.Token = uuid.Must(uuid.NewV4()).String()
	twoFa.ValidUntil = validUntil
	twoFa.ResendValidUntil = ResendValidUntil

	if err := rs.Store.CreateOrUpdateTwoFa(twoFa); err != nil {
		log().Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	message := fmt.Sprintf("Your Pathao Pay Login OTP is : %s", vc.Raw())
	err = vc.SendSMS(acc.Mobile, message)
	if err != nil {
		log().Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	mfa.SetChallengeHeader(w, twoFa)
	render.Respond(w, r, &resendOTPResponse{
		Resend: true,
	})
	return
}
