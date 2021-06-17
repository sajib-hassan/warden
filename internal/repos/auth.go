package repos

import (
	"time"

	"github.com/go-pg/pg"

	usingpin2 "github.com/sajib-hassan/warden/internal/auth/usingpin"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// AuthStore implements database operations for account PIN based authentication.
type AuthStore struct {
	db *pg.DB
}

// NewAuthStore return an AuthStore.
func NewAuthStore(db *pg.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

// GetUser returns an account by ID.
func (s *AuthStore) GetUser(id int) (*usingpin2.User, error) {
	a := usingpin2.User{ID: id}
	err := s.db.Model(&a).
		Column("user.*").
		Where("id = ?id").
		First()
	return &a, err
}

// GetUserByMobile returns an account by mobile.
func (s *AuthStore) GetUserByMobile(m string) (*usingpin2.User, error) {
	a := usingpin2.User{Mobile: m}
	err := s.db.Model(&a).
		Column("id", "active", "mobile", "name", "pin").
		Where("mobile = ?mobile").
		First()
	return &a, err
}

// UpdateUser upates account data related to PIN based authentication .
func (s *AuthStore) UpdateUser(a *usingpin2.User) error {
	_, err := s.db.Model(a).
		Column("last_login").
		WherePK().
		Update()
	return err
}

// GetToken returns refresh token by token identifier.
func (s *AuthStore) GetToken(t string) (*jwt.Token, error) {
	token := jwt.Token{Token: t}
	err := s.db.Model(&token).
		Where("token = ?token").
		First()

	return &token, err
}

// CreateOrUpdateToken creates or updates an existing refresh token.
func (s *AuthStore) CreateOrUpdateToken(t *jwt.Token) error {
	var err error
	if t.ID == 0 {
		err = s.db.Insert(t)
	} else {
		err = s.db.Update(t)
	}
	return err
}

// DeleteToken deletes a refresh token.
func (s *AuthStore) DeleteToken(t *jwt.Token) error {
	err := s.db.Delete(t)
	return err
}

// PurgeExpiredToken deletes expired refresh token.
func (s *AuthStore) PurgeExpiredToken() error {
	_, err := s.db.Model(&jwt.Token{}).
		Where("expiry < ?", time.Now()).
		Delete()

	return err
}
