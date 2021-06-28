package repos

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sajib-hassan/warden/internal/auth/usingpin"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
	"github.com/sajib-hassan/warden/pkg/auth/mfa"
)

// AuthStore implements database operations for user PIN based authentication.
type AuthStore struct {
}

// NewAuthStore return an AuthStore.
func NewAuthStore() *AuthStore {
	return &AuthStore{}
}

// GetUser returns an user by ID.
func (s *AuthStore) GetUser(id string) (*usingpin.User, error) {
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

// GetUserByMobile returns an user by mobile.
func (s *AuthStore) GetUserByMobile(m string) (*usingpin.User, error) {
	u := &usingpin.User{}
	err := mgm.Coll(u).First(bson.M{"mobile": m}, u)

	if err != nil {
		//if err == mongo.ErrNoDocuments {
		//	return nil, nil
		//}
		return nil, err
	}
	return u, nil
}

// UpdateUser updates user data related to PIN based authentication .
func (s *AuthStore) UpdateUser(u *usingpin.User) error {
	return mgm.Coll(u).Update(u)
}

// GetToken returns refresh token by token identifier.
func (s *AuthStore) GetToken(token string) (*jwt.Token, error) {
	t := &jwt.Token{}
	err := mgm.Coll(t).First(bson.M{"token": token}, t)
	if err != nil {
		//if err == mongo.ErrNoDocuments {
		//	return nil, nil
		//}
		return nil, err
	}
	return t, nil
}

// CreateOrUpdateToken creates or updates an existing refresh token.
func (s *AuthStore) CreateOrUpdateToken(t *jwt.Token) error {
	var err error
	if t.ID.IsZero() {
		err = mgm.Coll(t).Create(t)
	} else {
		err = mgm.Coll(t).Update(t)
	}
	return err
}

// DeleteToken deletes a refresh token.
func (s *AuthStore) DeleteToken(t *jwt.Token) error {
	return mgm.Coll(t).Delete(t)
}

// PurgeExpiredToken deletes expired refresh token.
func (s *AuthStore) PurgeExpiredToken() error {
	_, err := mgm.Coll(&jwt.Token{}).
		DeleteMany(mgm.Ctx(), bson.M{"expiry": bson.M{"$lt": time.Now()}})
	return err
}

func (s AuthStore) GetTrustedDevice(userId string, identifier string) (*usingpin.Device, error) {
	d := &usingpin.Device{}
	err := mgm.Coll(d).First(bson.M{"user_id": userId, "identifier": identifier, "is_authorized": true}, d)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

func (s AuthStore) RegisterAsTrustedDevice(d *usingpin.Device) error {
	return mgm.Coll(d).Create(d)
}

func (s AuthStore) CreateOrUpdateTwoFa(t *mfa.TwoFa) error {
	var err error
	if t.ID.IsZero() {
		err = mgm.Coll(t).Create(t)
	} else {
		err = mgm.Coll(t).Update(t)
	}
	return err
}
