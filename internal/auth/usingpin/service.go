package usingpin

import (
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/auth/mfa"
	"github.com/sajib-hassan/warden/pkg/notifier"
)

func SentLoginOTP(as AuthStorer, u *User) (string, error) {
	totp, err := GetTOTP(as, u)
	if err != nil {
		return "", err
	}

	challenge, err := totp.NowWithExpiration()
	if err != nil {
		return "", err
	}

	// Create OTP Object into DB
	token := uuid.Must(uuid.NewV4()).String()
	twoFa := &mfa.TwoFa{
		ServiceId:        u.ID.Hex(),
		ServiceType:      reflect.TypeOf(u).String(),
		Channel:          mfa.SMS,
		Mobile:           u.Mobile,
		Challenge:        challenge,
		Token:            token,
		MaxAllowedRetry:  viper.GetInt("AUTH_LOGIN_OTP_MAX_ALLOWED_RETRY"),
		MaxAllowedResend: viper.GetInt("AUTH_LOGIN_OTP_MAX_ALLOWED_RESEND"),
	}

	if err := as.CreateOrUpdateTwoFa(twoFa); err != nil {
		return "", err
	}

	smsClient, err := notifier.NewSMSClient()
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("Your Pathao Pay Login OTP is : %s", challenge)
	err = smsClient.Send(u.Mobile, message)

	if err != nil {
		return "", err
	}

	return token, nil

}

func GetTOTP(as AuthStorer, u *User) (*mfa.TOTP, error) {
	totp := mfa.NewTOTP(viper.GetUint("AUTH_LOGIN_OTP_TIMEOUT"))
	if u.OtpKey == "" {
		err := totp.NewKey(u.ID.Hex())
		if err != nil {
			return nil, err
		}

		u.OtpKey = totp.Key.String()
		err = as.UpdateUser(u)
		if err != nil {
			return nil, err
		}
	} else {
		err := totp.NewKeyFromURL(u.OtpKey)
		if err != nil {
			return nil, err
		}
	}
	return totp, nil
}
