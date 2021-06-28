package jwt

import (
	"time"

	"github.com/kamva/mgm/v3"
)

// Token holds refresh jwt information.
type Token struct {
	mgm.DefaultModel `bson:",inline"`
	UserID           string `json:"-" bson:"user_id"`

	Token      string    `json:"-" bson:"token"`
	Expiry     time.Time `json:"-" bson:"expiry"`
	Mobile     bool      `json:"mobile" bson:"mobile"`
	Identifier string    `json:"identifier,omitempty" bson:"identifier"`
}

// Claims returns the token claims to be signed
func (t *Token) Claims() RefreshClaims {
	return RefreshClaims{
		ID:    t.ID.Hex(),
		Token: t.Token,
	}
}
