package mfa

import "errors"

var (
	ErrChallengeTokenUnauthorized = errors.New("challenge unauthorized")
	ErrChallengeTokenExpired      = errors.New("challenge expired")
	ErrInvalidChallengeToken      = errors.New("invalid challenge")
	ErrResendTimeExpired          = errors.New("resend time is expired")
	ErrMaxResendExceed            = errors.New("resend limit is exceeded")
	ErrMaxRetryExceed             = errors.New("retry limit is exceeded")
	ErrChallengeCodeMissMatched   = errors.New("challenge not matched")
)
