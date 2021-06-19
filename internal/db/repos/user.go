package repos

import (
	"context"
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sajib-hassan/warden/internal/auth/usingpin"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// UserStore implements database operations for account management by user.
type UserStore struct {
	Coll *mgm.Collection
}

// NewUserStore returns an UserStore.
func NewUserStore() *UserStore {
	us := &UserStore{
		mgm.Coll(&usingpin.User{}),
	}
	us.EnsureIndices()
	return us
}

func (s *UserStore) EnsureIndices() error {
	log.Println("Starting EnsureIndices")
	_, err := s.Coll.Indexes().CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "mobile", Value: 1}},
				Options: options.Index().SetUnique(true).SetName("unique_mobile")},
		})
	log.Println("Completed EnsureIndices", err)
	return err
}

// Get an account by ID.
func (s *UserStore) Get(id int) (*usingpin.User, error) {
	a := usingpin.User{}
	//a := usingpin.User{ID: id}
	//err := s.db.Model(&a).
	//	Where("user.id = ?id").
	//	Column("user.*", "Token").
	//	First()
	//return &a, err
	return &a, nil
}

// Create a user
func (s *UserStore) Create(u *usingpin.User) error {
	return mgm.Coll(u).Create(u)
}

// Update an account.
func (s *UserStore) Update(a *usingpin.User) error {
	//_, err := s.db.Model(a).
	//	Column("mobile", "name").
	//	WherePK().
	//	Update()
	//return err
	return nil
}

// Delete an account.
func (s *UserStore) Delete(a *usingpin.User) error {
	//err := s.db.RunInTransaction(func(tx *pg.Tx) error {
	//	if _, err := tx.Model(&jwt.Token{}).
	//		Where("user_id = ?", a.ID).
	//		Delete(); err != nil {
	//		return err
	//	}
	//	if _, err := tx.Model(&models.Profile{}).
	//		Where("user_id = ?", a.ID).
	//		Delete(); err != nil {
	//		return err
	//	}
	//	return tx.Delete(a)
	//})
	//return err
	return nil
}

// UpdateToken updates a jwt refresh token.
func (s *UserStore) UpdateToken(t *jwt.Token) error {
	//_, err := s.db.Model(t).
	//	Column("identifier").
	//	WherePK().
	//	Update()
	//return err
	return nil
}

// DeleteToken deletes a jwt refresh token.
func (s *UserStore) DeleteToken(t *jwt.Token) error {
	//err := s.db.Delete(t)
	//return err
	return nil
}
