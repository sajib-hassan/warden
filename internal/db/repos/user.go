package repos

import (
	"github.com/kamva/mgm/v3"

	"github.com/sajib-hassan/warden/internal/auth/usingpin"
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
