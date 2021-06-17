package usingpin

import (
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/go-pg/pg/orm"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/auth/encryptor"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// User represents an authenticated application user
type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	LastLogin time.Time `json:"last_login,omitempty"`

	Mobile string   `json:"mobile"`
	Pin    string   `json:"pin"`
	Name   string   `json:"name"`
	Active bool     `sql:",notnull" json:"active"`
	Roles  []string `pg:",array" json:"roles,omitempty"`

	Token []jwt.Token `json:"token,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (u *User) BeforeInsert(db orm.DB) error {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
		u.UpdatedAt = now
	}

	err := u.validatePin()
	if err != nil {
		return err
	}
	return u.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (u *User) BeforeUpdate(db orm.DB) error {
	u.UpdatedAt = time.Now()
	return u.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (u *User) BeforeDelete(db orm.DB) error {
	return nil
}

func (u User) validatePin() error {
	loginPinLength := viper.GetInt("auth_login_pin_length")

	return validation.ValidateStruct(u,
		validation.Field(&u.Pin, validation.Required, is.Digit, validation.Length(loginPinLength, loginPinLength)),
	)
}

// Validate validates User struct and returns validation errors.
func (u *User) Validate() error {
	u.Mobile = strings.TrimSpace(u.Mobile)
	u.Mobile = strings.ToLower(u.Mobile)
	u.Name = strings.TrimSpace(u.Name)

	return validation.ValidateStruct(u,
		validation.Field(&u.Mobile, validation.Required, validation.By(CheckValidBDMobileNumber)),
		validation.Field(&u.Name, validation.Required, is.ASCII),
	)
}

// CanLogin returns true if user is allowed to login.
func (u *User) CanLogin() bool {
	return u.Active
}

func (u *User) isPinMatched(pin string) (bool, error) {
	return encryptor.ComparePasswordAndHash(pin, u.Pin)
}

// Claims returns the account's claims to be signed
func (u *User) Claims() jwt.AppClaims {
	return jwt.AppClaims{
		ID:    u.ID,
		Sub:   u.Name,
		Roles: u.Roles,
	}
}
