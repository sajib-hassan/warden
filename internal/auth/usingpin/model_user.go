package usingpin

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kamva/mgm/v3"
	"github.com/spf13/viper"

	"github.com/sajib-hassan/warden/pkg/auth/encryptor"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/validator"
)

// User represents an authenticated application user
type User struct {
	mgm.DefaultModel `bson:",inline"`

	Mobile string `json:"mobile" bson:"mobile"`
	Pin    string `json:"pin" bson:"pin"`
	Name   string `json:"name" bson:"name"`
	Active bool   `json:"active" bson:"active"`

	Secret    string    `json:"secret,omitempty" bson:"secret"`
	LastLogin time.Time `json:"last_login,omitempty" bson:"last_login"`

	Roles []string    `json:"roles,omitempty" bson:"roles"`
	Token []jwt.Token `json:"token,omitempty"`
}

// Creating hook executed before database insert operation.
//func (u *User) Creating() error {
//	// Call the DefaultModel Creating hook
//	if err := u.DefaultModel.Creating(); err != nil {
//		return err
//	}
//
//	err := u.validatePin()
//	if err != nil {
//		return err
//	}
//	return u.Validate()
//}

// Saving hook executed before database update operation.
//func (u *User) Saving() error {
//	return u.Validate()
//}

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
		validation.Field(&u.Mobile, validation.Required, validation.By(validator.CheckValidBDMobileNumber)),
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

// Claims returns the user's claims to be signed
func (u *User) Claims() jwt.AppClaims {
	return jwt.AppClaims{
		ID:    u.ID.Hex(),
		Sub:   u.Name,
		Roles: u.Roles,
	}
}
