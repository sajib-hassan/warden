package app

import (
	"github.com/sajib-hassan/warden/internal/auth/usingpin"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// UserStore defines database operations for user.
type UserStore interface {
	Get(id string) (*usingpin.User, error)
	Update(*usingpin.User) error
	Delete(*usingpin.User) error
	UpdateToken(*jwt.Token) error
	DeleteToken(*jwt.Token) error
}
