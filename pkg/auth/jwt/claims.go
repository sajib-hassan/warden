package jwt

import (
	"encoding/json"
	"errors"
)

type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

type MapClaims map[string]interface{}

// AppClaims represent the claims parsed from JWT access token.
type AppClaims struct {
	ID    string   `json:"id,omitempty"`
	Sub   string   `json:"sub,omitempty"`
	Token string   `json:"token,omitempty"`
	Roles []string `json:"roles,omitempty"`
	StandardClaims
}

// ParseClaims parses JWT claims into AppClaims.
func (c *AppClaims) ParseClaims(claims MapClaims) error {
	id, ok := claims["id"]
	if !ok {
		return errors.New("could not parse claim id")
	}
	c.ID = id.(string)

	sub, ok := claims["sub"]
	if !ok {
		return errors.New("could not parse claim sub")
	}
	c.Sub = sub.(string)

	token, ok := claims["token"]
	if !ok {
		return errors.New("could not parse access token")
	}
	c.Token = token.(string)

	rl, ok := claims["roles"]
	if !ok {
		return errors.New("could not parse claims roles")
	}

	var roles []string
	if rl != nil {
		for _, v := range rl.([]interface{}) {
			roles = append(roles, v.(string))
		}
	}
	c.Roles = roles

	return nil
}

func (c *AppClaims) AsMap() (m map[string]interface{}) {
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &m)
	return
}

// RefreshClaims represents the claims parsed from JWT refresh token.
type RefreshClaims struct {
	ID    string `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	StandardClaims
}

// ParseClaims parses the JWT claims into RefreshClaims.
func (c *RefreshClaims) ParseClaims(claims MapClaims) error {
	token, ok := claims["token"]
	if !ok {
		return errors.New("could not parse claim token")
	}
	c.Token = token.(string)
	return nil
}

func (c *RefreshClaims) AsMap() (m map[string]interface{}) {
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &m)
	return
}
