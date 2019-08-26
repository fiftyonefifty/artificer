package models

import "time"

type Token struct {
	ClientID         string        `json:clientID`
	UserID           string        `json:userID`
	RedirectURI      string        `json:redirectURI`
	Scope            string        `json:scope`
	Code             string        `json:code`
	CodeCreatedAt    time.Time     `json:codeCreatedAt`
	CodeExpiresIn    time.Duration `json:codeExpiresIn`
	Access           string        `json:access`
	AccessCreatedAt  time.Time     `json:accessCreatedAt`
	AccessExpiresIn  time.Duration `json:accessExpiresIn`
	Refresh          string        `json:refresh`
	RefreshCreatedAt time.Time     `json:refreshCreatedAt`
	RefreshExpiresIn time.Duration `json:refreshExpiresIn`
}
