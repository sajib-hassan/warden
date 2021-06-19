package repos

import (
	"errors"
	"net/url"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"

	"github.com/sajib-hassan/warden/internal/auth/usingpin"
	models2 "github.com/sajib-hassan/warden/internal/db/models"
	"github.com/sajib-hassan/warden/pkg/auth/jwt"
)

var (
	// ErrUniqueEmailConstraint provides error message for already registered email address.
	ErrUniqueEmailConstraint = errors.New("email already registered")
	// ErrBadParams could not parse params to filter
	ErrBadParams = errors.New("bad parameters")
)

// AdmAccountStore implements database operations for account management by admin.
type AdmAccountStore struct {
	db *pg.DB
}

// NewAdmAccountStore returns an UserStore.
func NewAdmAccountStore(db *pg.DB) *AdmAccountStore {
	return &AdmAccountStore{
		db: db,
	}
}

// AccountFilter provides pagination and filtering options on accounts.
type AccountFilter struct {
	Pager  *urlvalues.Pager
	Filter *urlvalues.Filter
	Order  []string
}

// NewAccountFilter returns an AccountFilter with options parsed from request url values.
func NewAccountFilter(params interface{}) (*AccountFilter, error) {
	v, ok := params.(url.Values)
	if !ok {
		return nil, ErrBadParams
	}
	p := urlvalues.Values(v)
	f := &AccountFilter{
		Pager:  urlvalues.NewPager(p),
		Filter: urlvalues.NewFilter(p),
		Order:  p["order"],
	}
	return f, nil
}

// Apply applies an AccountFilter on an orm.Query.
func (f *AccountFilter) Apply(q *orm.Query) (*orm.Query, error) {
	q = q.Apply(f.Pager.Pagination)
	q = q.Apply(f.Filter.Filters)
	q = q.Order(f.Order...)
	return q, nil
}

// List applies a filter and returns paginated array of matching results and total count.
func (s *AdmAccountStore) List(f *AccountFilter) ([]usingpin.User, int, error) {
	a := []usingpin.User{}
	count, err := s.db.Model(&a).
		Apply(f.Apply).
		SelectAndCount()
	if err != nil {
		return nil, 0, err
	}
	return a, count, nil
}

// Create creates a new account.
func (s *AdmAccountStore) Create(a *usingpin.User) error {
	count, _ := s.db.Model(a).
		Where("email = ?email").
		Count()

	if count != 0 {
		return ErrUniqueEmailConstraint
	}

	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		err := tx.Insert(a)
		if err != nil {
			return err
		}
		p := &models2.Profile{
			//UserID: a.ID,
		}
		return tx.Insert(p)
	})

	return err
}

// Get account by ID.
func (s *AdmAccountStore) Get(id int) (*usingpin.User, error) {
	a := usingpin.User{}
	//a := usingpin.User{ID: id}
	err := s.db.Select(&a)
	return &a, err
}

// Update account.
func (s *AdmAccountStore) Update(a *usingpin.User) error {
	err := s.db.Update(a)
	return err
}

// Delete account.
func (s *AdmAccountStore) Delete(a *usingpin.User) error {
	err := s.db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.Model(&jwt.Token{}).
			Where("user_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(&models2.Profile{}).
			Where("user_id = ?", a.ID).
			Delete(); err != nil {
			return err
		}
		return tx.Delete(a)
	})
	return err
}
