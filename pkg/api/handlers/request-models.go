package handlers

type (
	TokenRequest struct {
		ClientID     string `json:"client_id" form:"client_id" query:"client_id"`
		ClientSecret string `json:"client_secret" form:"client_secret" query:"client_secret"`
		GrantType    string `json:"grant_type" form:"grant_type" query:"grant_type"`
	}
	ClientCredentialsRequest struct {
		TokenRequest
		Scope string `json:"scope" form:"scope" query:"scope"`
	}

	ArbitraryNoSubjectRequest struct {
		TokenRequest
		ArbitraryClaims     string `json:"arbitrary_claims" form:"arbitrary_claims" query:"arbitrary_claims"`
		ArbitraryAmrs       string `json:"arbitrary_amrs" form:"arbitrary_amrs" query:"arbitrary_amrs"`
		ArbitraryAudiences  string `json:"arbitrary_audiences" form:"arbitrary_audiences" query:"arbitrary_audiences"`
		AccessTokenLifetime int    `json:"access_token_lifetime" form:"access_token_lifetime" query:"access_token_lifetime"`
	}

	ClientCredentialsResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
)
