package mfa

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kamva/mgm/v3"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/validator"
)

const SMS = "sms"

const (
	STATUS_CREATED  = "created"
	STATUS_SENT     = "sent"
	STATUS_RESENT   = "resent"
	STATUS_VERIFIED = "verified"
	STATUS_EXPIRE   = "expire"
)

type TwoFa struct {
	mgm.DefaultModel   `bson:",inline"`
	ServiceId          string `json:"service_id" bson:"service_id"`
	ServiceType        string `json:"service_type" bson:"service_type"`
	Channel            string `json:"channel" bson:"channel"` // Default is 'sms'
	Mobile             string `json:"mobile" bson:"mobile"`
	Challenge          string `json:"challenge" bson:"challenge"`
	Token              string `json:"token" bson:"token"`
	ValidUntil         int64  `json:"valid_until" bson:"valid_until"`
	RetryCount         int    `json:"retry_count" bson:"retry_count"`
	MaxAllowedRetry    int    `json:"max_allowed_retry" bson:"max_allowed_retry"`
	ResendCount        int    `json:"resend_count" bson:"resend_count"`
	MaxAllowedResend   int    `json:"max_allowed_resend" bson:"max_allowed_resend"`
	ChallengeSignature string `json:"challenge_signature" bson:"challenge_signature"`
	IdentitySignature  string `json:"identity_signature" bson:"identity_signature"`
	UsedFor            string `json:"used_for" bson:"used_for"`
	Status             string `json:"status" bson:"status"`
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

	if t.Status == "" {
		t.Status = STATUS_CREATED
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
		validation.Field(&t.MaxAllowedRetry, validation.Required),
		validation.Field(&t.MaxAllowedResend, validation.Required),
		validation.Field(&t.UsedFor, validation.Required, is.ASCII),
		validation.Field(&t.Status, validation.Required, is.ASCII),
	)
}
