// Package models contains application specific entities.
package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/go-pg/pg/orm"
)

// Profile holds specific application settings linked to an User.
type Profile struct {
	ID        int       `json:"-"`
	UserID    int       `json:"-"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Nid         string `json:"nid,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}

// BeforeInsert hook executed before database insert operation.
func (p *Profile) BeforeInsert(db orm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook executed before database update operation.
func (p *Profile) BeforeUpdate(db orm.DB) error {
	p.UpdatedAt = time.Now()
	return p.Validate()
}

// Validate validates Profile struct and returns validation errors.
func (p *Profile) Validate() error {

	return validation.ValidateStruct(p,
		validation.Field(&p.Nid, validation.Required, validation.Length(10, 17), is.Digit),
		validation.Field(&p.DateOfBirth, validation.Required),
	)
}
