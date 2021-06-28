// Package models contains application specific entities.
package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kamva/mgm/v3"
)

// Profile holds specific application settings linked to an User.
type Profile struct {
	mgm.DefaultModel `bson:",inline"`
	UserID           string `json:"-"`

	Nid         string    `json:"nid,omitempty" bson:"nid"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty" bson:"date_of_birth"`
}

// Saving hook executed before database update operation.
func (p *Profile) Saving() error {
	return p.Validate()
}

// Validate validates Profile struct and returns validation errors.
func (p *Profile) Validate() error {

	return validation.ValidateStruct(p,
		validation.Field(&p.Nid, validation.Required, validation.Length(10, 17), is.Digit),
		validation.Field(&p.DateOfBirth, validation.Required),
	)
}
