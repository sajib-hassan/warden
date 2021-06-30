package mfa

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sajib-hassan/warden/pkg/helpmate"
	"github.com/sajib-hassan/warden/pkg/notifier"
)

type VerificationCode struct {
	code string
}

func NewVerificationCode() *VerificationCode {
	return &VerificationCode{}
}

func (c *VerificationCode) randCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rndCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return rndCode
}

// Raw 6-digit random verification code
func (c *VerificationCode) Raw() string {
	if c.code == "" {
		c.code = c.randCode()
	}
	return c.code
}

func (c *VerificationCode) Hashed(secret string) string {
	if c.code == "" {
		c.code = c.randCode()
	}
	return helpmate.SHA256HMACHash(c.code, secret)
}

func (c *VerificationCode) IsValidChallengeCode(code string, secret string, hash string) bool {
	return helpmate.SHA256HMACHash(code, secret) == hash
}

func (c *VerificationCode) SendSMS(to, message string) error {
	smsClient, err := notifier.NewSMSClient()
	if err != nil {
		return err
	}
	err = smsClient.Send(to, message)
	return err
}
