package app

import (
	"github.com/sajib-hassan/warden/internal/db/models"
	"github.com/sajib-hassan/warden/pkg/auth/authorize"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// UserStore defines database operations for user.
type UserStore interface {
	Get(id string) (*authorize.User, error)
	Update(*authorize.User) error
	Delete(*authorize.User) error
	UpdateToken(*jwt.Token) error
	DeleteToken(*jwt.Token) error
}

// ProfileStore defines database operations for a profile.
type ProfileStore interface {
	Get(userID string) (*models.Profile, error)
	Update(p *models.Profile) error
}
