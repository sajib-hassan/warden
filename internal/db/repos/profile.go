package repos

import (
	"github.com/go-pg/pg"

	models2 "github.com/sajib-hassan/warden/internal/db/models"
)

// ProfileStore implements database operations for profile management.
type ProfileStore struct {
	db *pg.DB
}

// NewProfileStore returns a ProfileStore implementation.
func NewProfileStore() *ProfileStore {
	db := &pg.DB{}
	return &ProfileStore{
		db: db,
	}
}

// Get gets an profile by account ID.
func (s *ProfileStore) Get(accountID int) (*models2.Profile, error) {
	p := models2.Profile{UserID: accountID}
	_, err := s.db.Model(&p).
		Where("user_id = ?", accountID).
		SelectOrInsert()

	return &p, err
}

// Update updates profile.
func (s *ProfileStore) Update(p *models2.Profile) error {
	err := s.db.Update(p)
	return err
}
