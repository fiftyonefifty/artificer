package models

type TokenUsage int

const (
	ReUse TokenUsage = iota
	OneTimeOnly
)

type TokenExpiration int

const (
	Sliding TokenExpiration = iota
	Absolute
)

type AccessTokenType int

const (
	Jwt AccessTokenType = iota
	Reference
)

type Claim struct {
	Type  string `json:type`
	Value string `json:value`
}

type Secret struct {
	Value      string `json:value`
	Expiration int    `json:expiration`
}

type Client struct {
	Enabled                          bool            `json:enabled`
	ClientID                         string          `json:clientID`
	ClientName                       string          `json:clientName`
	Description                      string          `json:description`
	Namespace                        string          `json:namespace`
	RequireRefreshClientSecret       bool            `json:requireRefreshClientSecret`
	AllowOfflineAccess               bool            `json:allowOfflineAccess`
	AccessTokenLifetime              int             `json:accessTokenLifetime`
	AbsoluteRefreshTokenLifetime     int             `json:absoluteRefreshTokenLifetime`
	SlidingRefreshTokenLifetime      int             `json:slidingRefreshTokenLifetime`
	UpdateAccessTokenClaimsOnRefresh bool            `json:updateAccessTokenClaimsOnRefresh`
	RefreshTokenUsage                TokenUsage      `json:refreshTokenUsage`
	RefreshTokenExpiration           TokenExpiration `json:refreshTokenExpiration`
	AccessTokenType                  AccessTokenType `json:accessTokenType`
	IncludeJwtId                     bool            `json:includeJwtId`
	Claims                           []Claim         `json:claims`
	AlwaysSendClientClaims           bool            `json:alwaysSendClientClaims`
	AlwaysIncludeUserClaimsInIdToken bool            `json:alwaysIncludeUserClaimsInIdToken`
	AllowedScopes                    []string        `json:allowedScopes`
	AllowedGrantTypes                []string        `json:allowedGrantTypes`
	RequireClientSecret              bool            `json:requireClientSecret`
	ClientSecrets                    []Secret        `json:clientSecrets`
}
