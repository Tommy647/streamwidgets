package twitch

import (
	"encoding/json"
	"time"
)

// clientToken response
type clientToken struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken *string       `json:"refresh_token"`
	ExpiresIn    time.Duration `json:"expires_in"`
	Scope        []string      `json:"scope"`
	TokenType    string        `json:"token_type"`
	Error        error         `json:"_,omitempty"`
}

// UnmarshalJSON custom unmarshaller to get the correct expiry time
func (c *clientToken) UnmarshalJSON(bytes []byte) error {
	type alias clientToken
	var a alias
	if err := json.Unmarshal(bytes, &a); err != nil {
		return err
	}
	*c = (clientToken)(a)
	c.ExpiresIn = a.ExpiresIn * time.Second
	return nil
}

func (c clientToken) MarshalJSON() ([]byte, error) {
	type alias clientToken
	var a = (alias)(c)
	a.ExpiresIn = c.ExpiresIn / time.Second
	return json.Marshal(a)
}
