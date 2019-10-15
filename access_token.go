package gotwilio

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AccessToken holds everything required to
// create authenticate the client SDKs.
// See https://www.twilio.com/docs/iam/access-tokens
// for further details.
type AccessToken struct {
	AccountSid   string
	APIKeySid    string
	APIKeySecret string

	NotBefore time.Time
	ExpiresAt time.Time
	Grants    []Grant
	Identity  string
}

// Grant is a perimssion given to the Access Token.
// Types include Chat, Video etc.
type Grant interface {
	GrantName() string
}

// VideoGrant is the permission to use the Video API
// which can be given to an Access Token.
type VideoGrant struct {
	Room string `json:"room,omitempty"`
}

func (g *VideoGrant) GrantName() string {
	return "video"
}

// NewAccessToken creates a new Access Token which
// can be used to authenticate Twilio Client SDKs
// for a short period of time.
func (twilio *Twilio) NewAccessToken() *AccessToken {
	return &AccessToken{
		AccountSid:   twilio.AccountSid,
		APIKeySid:    twilio.APIKeySid,
		APIKeySecret: twilio.APIKeySecret,
	}
}

// AddGrant adds a given Grant to the Access Token.
func (a *AccessToken) AddGrant(grant Grant) *AccessToken {
	a.Grants = append(a.Grants, grant)
	return a
}

// ToJWT creates a JSON Web Token from the Access Token
// to use in the Client SDKs.
// See https://en.wikipedia.org/wiki/JSON_Web_Token
// for the standard format.
func (a *AccessToken) ToJWT() (string, error) {
	claims := &twilioClaims{
		jwt.StandardClaims{
			Id:        a.APIKeySid + fmt.Sprintf("-%d", time.Now().UnixNano()),
			Issuer:    a.APIKeySid,
			Subject:   a.AccountSid,
			NotBefore: a.NotBefore.Unix(),
			ExpiresAt: a.ExpiresAt.Unix(),
		},
		&grantsClaim{
			Identity: a.Identity,
			Grants:   a.Grants,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header = map[string]interface{}{
		"typ": "JWT",
		"alg": "HS256",
		"cty": "twilio-fpa;v=1",
	}

	ss, err := token.SignedString([]byte(a.APIKeySecret))

	return ss, err
}

// Private helpers to construct the JWT.

type twilioClaims struct {
	jwt.StandardClaims
	Grants *grantsClaim `json:"grants"`
}

type grantsClaim struct {
	Identity string
	Grants   []Grant
}

func (g *grantsClaim) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})
	data["identity"] = g.Identity
	for _, grant := range g.Grants {
		data[grant.GrantName()] = grant
	}
	return json.Marshal(data)
}
