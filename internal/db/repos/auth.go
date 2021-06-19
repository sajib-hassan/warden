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

// AuthStore implements database operations for account PIN based authentication.
type AuthStore struct {
	Coll *mgm.Collection
}

// NewAuthStore return an AuthStore.
func NewAuthStore() *AuthStore {
	as := &AuthStore{
		mgm.Coll(&usingpin.User{}),
	}
	as.EnsureIndices()
	return as
}

func (s *AuthStore) EnsureIndices() error {
	log.Println("Starting EnsureIndices")
	_, err := s.Coll.Indexes().CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "mobile", Value: 1}},
				Options: options.Index().SetUnique(true)},
		})
	log.Println("Completed EnsureIndices", err)
	return err
}

// GetUser returns an account by ID.
func (s *AuthStore) GetUser(id int) (*usingpin.User, error) {
	a := usingpin.User{}
	//a := usingpin.User{ID: id}

	//err := s.db.Model(&a).
	//	Column("user.*").
	//	Where("id = ?id").
	//	First()
	//return &a, err
	return &a, nil
}

// GetUserByMobile returns an account by mobile.
func (s *AuthStore) GetUserByMobile(m string) (*usingpin.User, error) {
	user := &usingpin.User{}
	err := s.Coll.First(bson.M{"mobile": m}, user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// UpdateUser upates account data related to PIN based authentication .
func (s *AuthStore) UpdateUser(a *usingpin.User) error {
	//_, err := s.db.Model(a).
	//	Column("last_login").
	//	WherePK().
	//	Update()
	//return err
	return nil
}

// GetToken returns refresh token by token identifier.
func (s *AuthStore) GetToken(t string) (*jwt.Token, error) {
	token := jwt.Token{Token: t}
	//err := s.db.Model(&token).
	//	Where("token = ?token").
	//	First()
	//
	//return &token, err
	return &token, nil
}

// CreateOrUpdateToken creates or updates an existing refresh token.
func (s *AuthStore) CreateOrUpdateToken(t *jwt.Token) error {
	var err error
	//if t.ID == 0 {
	//	err = s.db.Insert(t)
	//} else {
	//	err = s.db.Update(t)
	//}
	return err
}

// DeleteToken deletes a refresh token.
func (s *AuthStore) DeleteToken(t *jwt.Token) error {
	//err := s.db.Delete(t)
	//return err
	return nil
}

// PurgeExpiredToken deletes expired refresh token.
func (s *AuthStore) PurgeExpiredToken() error {
	//_, err := s.db.Model(&jwt.Token{}).
	//	Where("expiry < ?", time.Now()).
	//	Delete()
	//
	//return err
	return nil
}
