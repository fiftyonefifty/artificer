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
)
