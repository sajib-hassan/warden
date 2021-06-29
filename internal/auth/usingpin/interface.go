package usingpin

import (
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/auth/mfa"
)

// AuthStorer defines database operations on users and tokens.
type AuthStorer interface {
	GetUser(id string) (*User, error)
	GetUserByMobile(mobile string) (*User, error)
	UpdateUser(a *User) error

	GetToken(token string) (*jwt.Token, error)
	CreateOrUpdateToken(t *jwt.Token) error
	DeleteToken(t *jwt.Token) error
	PurgeExpiredToken() error

	GetTrustedDevice(userId string, did string) (*Device, error)
	RegisterAsTrustedDevice(d *Device) error

	CreateOrUpdateTwoFa(t *mfa.TwoFa) error
}
