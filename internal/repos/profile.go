package repos

import (
	"github.com/go-pg/pg"

	"github.com/sajib-hassan/warden/internal/models"
)

// ProfileStore implements database operations for profile management.
type ProfileStore struct {
	db *pg.DB
}

// NewProfileStore returns a ProfileStore implementation.
func NewProfileStore(db *pg.DB) *ProfileStore {
	return &ProfileStore{
		db: db,
	}
}

// Get gets an profile by account ID.
func (s *ProfileStore) Get(accountID int) (*models.Profile, error) {
	p := models.Profile{UserID: accountID}
	_, err := s.db.Model(&p).
		Where("user_id = ?", accountID).
		SelectOrInsert()

	return &p, err
}

// Update updates profile.
func (s *ProfileStore) Update(p *models.Profile) error {
	err := s.db.Update(p)
	return err
}
