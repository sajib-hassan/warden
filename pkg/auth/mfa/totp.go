package mfa

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/viper"
)

type TOTP struct {
	Key    *otp.Key
	Period uint
}

func NewTOTP(p uint) *TOTP {
	if p == 0 {
		p = 30
	}
	return &TOTP{Period: p}
}

func (t *TOTP) NewKey(name string) error {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      viper.GetString("auth_totp_issuer"),
		AccountName: name,
	})
	if err != nil {
		return err
	}

	t.Key = key

	return nil
}

func (t *TOTP) NewKeyFromURL(url string) error {
	key, err := otp.NewKeyFromURL(url)
	if err != nil {
		return err
	}

	t.Key = key
	return err
}

func (t *TOTP) NowWithDefault() (string, error) {
	return totp.GenerateCode(t.Key.Secret(), time.Now().UTC())
}

func (t *TOTP) NowWithExpiration() (string, error) {
	return totp.GenerateCodeCustom(t.Key.Secret(), time.Now().UTC(), totp.ValidateOpts{
		Period:    t.Period,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
}

func (t *TOTP) ValidatePassCode(passcode string) (bool, error) {
	return totp.ValidateCustom(
		passcode,
		t.Key.Secret(),
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    t.Period,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
}
