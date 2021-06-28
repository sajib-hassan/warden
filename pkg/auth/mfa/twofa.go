package mfa

import (
	"github.com/kamva/mgm/v3"
)

const SMS = "sms"

type TwoFa struct {
	mgm.DefaultModel   `bson:",inline"`
	ServiceId          string `json:"service_id" bson:"service_id"`
	ServiceType        string `json:"service_type" bson:"service_type"`
	Channel            string `json:"channel" bson:"channel"` // Default is 'sms'
	Mobile             string `json:"mobile" bson:"mobile"`
	Challenge          string `json:"challenge" bson:"challenge"`
	Token              string `json:"token" bson:"token"`
	RetryCount         int    `json:"retry_count" bson:"retry_count"`
	MaxAllowedRetry    int    `json:"max_allowed_retry" bson:"max_allowed_retry"`
	ResendCount        int    `json:"resend_count" bson:"resend_count"`
	MaxAllowedResend   int    `json:"max_allowed_resend" bson:"max_allowed_resend"`
	ChallengeSignature string `json:"challenge_signature" bson:"challenge_signature"`
	IdentitySignature  string `json:"identity_signature" bson:"identity_signature"`
	UsedFor            string `json:"used_for" bson:"used_for"`
	Status             string `json:"status" bson:"status"`
}
