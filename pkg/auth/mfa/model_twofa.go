package mfa

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kamva/mgm/v3"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/validator"
)

const SMS = "sms"

const (
	StatusCreated = "created"
)

type TwoFa struct {
	mgm.DefaultModel   `bson:",inline"`
	ServiceId          string                 `json:"service_id" bson:"service_id"`
	ServiceType        string                 `json:"service_type" bson:"service_type"`
	Channel            string                 `json:"channel" bson:"channel"` // Default is 'sms'
	Mobile             string                 `json:"mobile" bson:"mobile"`
	Challenge          string                 `json:"challenge" bson:"challenge"`
	Token              string                 `json:"token" bson:"token"`
	ValidUntil         int64                  `json:"valid_until" bson:"valid_until"`
	RetryCount         int                    `json:"retry_count" bson:"retry_count"`
	MaxAllowedRetry    int                    `json:"max_allowed_retry" bson:"max_allowed_retry"`
	ResendCount        int                    `json:"resend_count" bson:"resend_count"`
	MaxAllowedResend   int                    `json:"max_allowed_resend" bson:"max_allowed_resend"`
	ResendValidUntil   int64                  `json:"resend_valid_until" bson:"resend_valid_until"`
	ChallengeSignature string                 `json:"challenge_signature,omitempty" bson:"challenge_signature"`
	IdentitySignature  string                 `json:"identity_signature,omitempty" bson:"identity_signature"`
	UsedFor            string                 `json:"used_for" bson:"used_for"`
	MetaData           map[string]interface{} `json:"meta_data,omitempty" bson:"meta_data"`
}

//Creating hook executed before database insert operation.
func (t *TwoFa) Creating() error {
	// Call the DefaultModel Creating hook
	if err := t.DefaultModel.Creating(); err != nil {
		return err
	}

	if t.MaxAllowedRetry == 0 {
		t.MaxAllowedRetry = viper.GetInt("AUTH_LOGIN_OTP_MAX_ALLOWED_RETRY")
	}

	if t.MaxAllowedResend == 0 {
		t.MaxAllowedResend = viper.GetInt("AUTH_LOGIN_OTP_MAX_ALLOWED_RESEND")
	}

	return t.Validate()
}

// Validate validates User struct and returns validation errors.
func (t *TwoFa) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.ServiceId, validation.Required, is.ASCII),
		validation.Field(&t.ServiceType, validation.Required, is.ASCII),
		validation.Field(&t.Mobile, validation.Required, validation.By(validator.CheckValidBDMobileNumber)),
		validation.Field(&t.Challenge, validation.Required, is.ASCII),
		validation.Field(&t.Token, validation.Required, is.UUIDv4),
		validation.Field(&t.ValidUntil, validation.Required),
		validation.Field(&t.ResendValidUntil, validation.Required),
		validation.Field(&t.MaxAllowedRetry, validation.Required),
		validation.Field(&t.MaxAllowedResend, validation.Required),
		validation.Field(&t.UsedFor, validation.Required, is.ASCII),
	)
}

func (t *TwoFa) IsExpired() bool {
	//return time.Now().UnixNano() > t.ValidUntil
	vu := time.Unix(t.ValidUntil, 0)
	return time.Now().UTC().After(vu)
}

func (t *TwoFa) IsResendExpired() bool {
	//return time.Now().UnixNano() > t.ResendValidUntil
	vu := time.Unix(t.ResendValidUntil, 0)
	return time.Now().UTC().After(vu)
}

func (t *TwoFa) IsExceedMaxRetry() bool {
	return t.RetryCount >= t.MaxAllowedRetry
}

func (t *TwoFa) IsExceedMaxResend() bool {
	return t.ResendCount >= t.MaxAllowedResend
}

func (t *TwoFa) IncreaseResendCount() {
	t.ResendCount = t.ResendCount + 1
}
