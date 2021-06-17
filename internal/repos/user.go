package repos

import (
	"github.com/go-pg/pg"

	usingpin2 "github.com/sajib-hassan/warden/internal/auth/usingpin"
	"github.com/sajib-hassan/warden/internal/models"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

// UserStore implements database operations for account management by user.
type UserStore struct {
	db *pg.DB
}

// NewAccountStore returns an UserStore.
func NewAccountStore(db *pg.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// Get an account by ID.
func (s *UserStore) Get(id int) (*usingpin2.User, error) {
	a := usingpin2.User{ID: id}
	err := s.db.Model(&a).
		Where("user.id = ?id").
		Column("user.*", "Token").
		First()
	return &a, err
}

// Update an account.
func (s *UserStore) Update(a *usingpin2.User) error {
	_, err := s.db.Model(a).
		Column("mobile", "name").
		WherePK().
		Update()
	return err
}

// Delete an account.
func (s *UserStore) Delete(a *usingpin2.User) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&jwt.Token{}).
			Where("user_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(&models.Profile{}).
			Where("user_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		return tx.Delete(a)
	})
	return err
}

// UpdateToken updates a jwt refresh token.
func (s *UserStore) UpdateToken(t *jwt.Token) error {
	_, err := s.db.Model(t).
		Column("identifier").
		WherePK().
		Update()
	return err
}

// DeleteToken deletes a jwt refresh token.
func (s *UserStore) DeleteToken(t *jwt.Token) error {
	err := s.db.Delete(t)
	return err
}
