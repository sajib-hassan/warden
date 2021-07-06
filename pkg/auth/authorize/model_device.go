package authorize

import "github.com/kamva/mgm/v3"

type Device struct {
	mgm.DefaultModel `bson:",inline"`

	UserID       string                 `json:"-" bson:"user_id"`
	Identifier   string                 `json:"identifier" bson:"identifier"`
	Name         string                 `json:"name" bson:"name"`
	IsAuthorized bool                   `json:"is_authorized" bson:"is_authorized"`
	Details      map[string]interface{} `json:"details,omitempty" bson:"details"`
}
