package repos

import (
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sajib-hassan/warden/pkg/auth/authorize"
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
func (s *AuthStore) GetUser(id string) (*authorize.User, error) {
	u := &authorize.User{}
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
func (s *AuthStore) GetUserByMobile(m string) (*authorize.User, error) {
	u := &authorize.User{}
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
func (s *AuthStore) UpdateUser(u *authorize.User) error {
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
		DeleteMany(mgm.Ctx(), bson.M{"expiry": bson.M{operator.Lt: time.Now()}})
	return err
}

func (s AuthStore) GetTrustedDevice(userId string, identifier string) (*authorize.Device, error) {
	d := &authorize.Device{}
	err := mgm.Coll(d).First(bson.M{"user_id": userId, "identifier": identifier, "is_authorized": true}, d)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

func (s AuthStore) RegisterAsTrustedDevice(d *authorize.Device) error {
	var err error
	if d.ID.IsZero() {
		err = mgm.Coll(d).Create(d)
	} else {
		err = mgm.Coll(d).Update(d)
	}
	return err
}

func (s AuthStore) GetTwoFa(token string, serviceType string, usedFor string) (*mfa.TwoFa, error) {
	t := &mfa.TwoFa{}
	err := mgm.Coll(t).First(bson.M{"token": token, "service_type": serviceType, "used_for": usedFor}, t)
	if err != nil {
		//if err == mongo.ErrNoDocuments {
		//	return nil, nil
		//}
		return nil, err
	}
	return t, nil
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

func (s *AuthStore) DeleteTwoFa(t *mfa.TwoFa) error {
	return mgm.Coll(t).Delete(t)
}

func (s *AuthStore) PurgeExpiredTwoFa() error {
	_, err := mgm.Coll(&mfa.TwoFa{}).
		DeleteMany(mgm.Ctx(), bson.M{"updated_at": bson.M{operator.Lt: time.Now().UTC().Add(time.Hour * -1)}})
	return err
}
