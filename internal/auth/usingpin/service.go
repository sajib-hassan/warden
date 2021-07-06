package usingpin

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"github.com/mssola/user_agent"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/auth/authorize"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/auth/mfa"
	"github.com/sajib-hassan/warden/pkg/helpmate"
)

type service struct {
	w  http.ResponseWriter
	r  *http.Request
	rs *Resource
}

func (s *service) performLogin(acc *authorize.User) {
	ua := user_agent.New(s.r.UserAgent())
	browser, _ := ua.Browser()

	token := &jwt.Token{
		Token:      uuid.Must(uuid.NewV4()).String(),
		Expiry:     time.Now().Add(s.rs.TokenAuth.JwtRefreshExpiry),
		UserID:     acc.ID.Hex(),
		Mobile:     ua.Mobile(),
		Identifier: fmt.Sprintf("%s on %s", browser, ua.OS()),
	}

	access, refresh, done := s.getTokens(acc, token)
	if !done {
		return
	}

	render.Respond(s.w, s.r, &loginResponse{
		Access:  access,
		Refresh: refresh,
		Slug:    acc.ID.Hex(),
		Name:    acc.Name,
		Mobile:  acc.Mobile,
	})
}

func (s *service) getTokens(acc *authorize.User, token *jwt.Token) (string, string, bool) {
	if err := s.rs.Store.CreateOrUpdateToken(token); err != nil {
		log().Error(err)
		render.Render(s.w, s.r, ErrInternalServerError)
		return "", "", false
	}

	access, refresh, err := s.rs.TokenAuth.GenTokenPair(acc.Claims(token), token.Claims())
	if err != nil {
		log().Error(err)
		render.Render(s.w, s.r, ErrInternalServerError)
		return "", "", false
	}

	acc.LastLogin = time.Now()
	if err := s.rs.Store.UpdateUser(acc); err != nil {
		log().Error(err)
		render.Render(s.w, s.r, ErrInternalServerError)
		return "", "", false
	}
	return access, refresh, true
}

func (s *service) sentLoginOTP(u *authorize.User, body *loginRequest) (*mfa.TwoFa, error) {
	vc := mfa.NewVerificationCode()

	if u.Secret == "" {
		u.Secret = helpmate.RandSecret(20)
		if err := s.rs.Store.UpdateUser(u); err != nil {
			return nil, err
		}
	}

	// Create OTP Object into DB
	validUntil := time.Now().UTC().Add(viper.GetDuration("AUTH_LOGIN_OTP_TIMEOUT")).Unix()
	ResendValidUntil := time.Now().UTC().Add(viper.GetDuration("AUTH_LOGIN_OTP_RESEND_TIMEOUT")).Unix()
	twoFa := &mfa.TwoFa{
		ServiceId:        u.ID.Hex(),
		ServiceType:      reflect.TypeOf(u).String(),
		Channel:          mfa.SMS,
		Mobile:           u.Mobile,
		Challenge:        vc.Hashed(u.Secret),
		Token:            uuid.Must(uuid.NewV4()).String(),
		ValidUntil:       validUntil,
		ResendValidUntil: ResendValidUntil,
		UsedFor:          "login",
		MetaData: map[string]interface{}{
			"device_id": body.DeviceId,
		},
	}

	if err := s.rs.Store.CreateOrUpdateTwoFa(twoFa); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Your Pathao Pay Login OTP is : %s", vc.Raw())
	err := vc.SendSMS(u.Mobile, message)
	if err != nil {
		return nil, err
	}

	return twoFa, nil
}

//func SentTOTP(as AuthStorer, u *User) (string, error) {
//	totp, err := GetTOTP(as, u)
//	if err != nil {
//		return "", err
//	}
//
//	challenge, err := totp.NowWithExpiration()
//	if err != nil {
//		return "", err
//	}
//
//	// Create OTP Object into DB
//	token := uuid.Must(uuid.NewV4()).String()
//	twoFa := &mfa.TwoFa{
//		ServiceId:        u.ID.Hex(),
//		ServiceType:      reflect.TypeOf(u).String(),
//		Channel:          mfa.SMS,
//		Mobile:           u.Mobile,
//		Challenge:        challenge,
//		Token:            token,
//		MaxAllowedRetry:  viper.GetInt("AUTH_LOGIN_OTP_MAX_ALLOWED_RETRY"),
//		MaxAllowedResend: viper.GetInt("AUTH_LOGIN_OTP_MAX_ALLOWED_RESEND"),
//	}
//
//	if err := as.CreateOrUpdateTwoFa(twoFa); err != nil {
//		return "", err
//	}
//
//	smsClient, err := notifier.NewSMSClient()
//	if err != nil {
//		return "", err
//	}
//	message := fmt.Sprintf("Your Pathao Pay Login OTP is : %s", challenge)
//	err = smsClient.Send(u.Mobile, message)
//
//	if err != nil {
//		return "", err
//	}
//
//	return token, nil
//
//}

//func GetTOTP(as AuthStorer, u *User) (*mfa.TOTP, error) {
//	totp := mfa.NewTOTP(viper.GetUint("AUTH_LOGIN_OTP_TIMEOUT"))
//	if u.Secret == "" {
//		err := totp.NewKey(u.ID.Hex())
//		if err != nil {
//			return nil, err
//		}
//
//		u.Secret = totp.Key.String()
//		err = as.UpdateUser(u)
//		if err != nil {
//			return nil, err
//		}
//	} else {
//		err := totp.NewKeyFromURL(u.Secret)
//		if err != nil {
//			return nil, err
//		}
//	}
//	return totp, nil
//}
