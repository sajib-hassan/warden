package repos

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sajib-hassan/warden/internal/auth/usingpin"
	"github.com/sajib-hassan/warden/internal/db/models"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// UserStore implements database operations for account management by user.
type UserStore struct {
}

// NewUserStore returns an UserStore.
func NewUserStore() *UserStore {
	return &UserStore{}
}

// Get an account by ID.
func (s *UserStore) Get(id string) (*usingpin.User, error) {
	u := &usingpin.User{}
	err := mgm.Coll(u).FindByID(id, u)

	if err != nil {
		//if err == mongo.ErrNoDocuments {
		//	return nil, nil
		//}
		return nil, err
	}
	return u, nil
}

// Create a user
func (s *UserStore) Create(u *usingpin.User) error {
	return mgm.Coll(u).Create(u)
}

// Update an account.
func (s *UserStore) Update(u *usingpin.User) error {
	return mgm.Coll(u).Update(u)
}

// Delete an account.
func (s *UserStore) Delete(u *usingpin.User) error {

	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {

		_, err := mgm.Coll(&jwt.Token{}).
			DeleteMany(sc, bson.M{"user_id": u.ID.Hex()})
		if err != nil {
			return err
		}

		_, err = mgm.Coll(&models.Profile{}).
			DeleteMany(sc, bson.M{"user_id": u.ID.Hex()})
		if err != nil {
			return err
		}

		err = mgm.Coll(u).DeleteWithCtx(sc, u)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})
}

// UpdateToken updates a jwt refresh token.
func (s *UserStore) UpdateToken(t *jwt.Token) error {
	return mgm.Coll(t).Update(t)
}

// DeleteToken deletes a jwt refresh token.
func (s *UserStore) DeleteToken(t *jwt.Token) error {
	return mgm.Coll(t).Delete(t)
}
