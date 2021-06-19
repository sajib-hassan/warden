package repos

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sajib-hassan/warden/internal/db/models"
)

// ProfileStore implements database operations for profile management.
type ProfileStore struct{}

// NewProfileStore returns a ProfileStore implementation.
func NewProfileStore() *ProfileStore {
	return &ProfileStore{}
}

// Get gets an profile by user ID.
func (s *ProfileStore) Get(userID string) (*models.Profile, error) {
	p := &models.Profile{UserID: userID}
	coll := mgm.Coll(p)

	err := coll.First(bson.M{"user_id": userID}, p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = coll.Create(p)
			if err == nil {
				return p, nil
			}
		}
		return nil, err
	}
	return p, nil
}

// Update updates profile.
func (s *ProfileStore) Update(p *models.Profile) error {
	return mgm.Coll(p).Update(p)
}
