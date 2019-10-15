package renderings

type WellKnownOpenidConfigurationResponse struct {
	Issuer                             string   `json:"issuer"`
	JwksURI                            string   `json:"jwks_uri"`
	AuthorizationEndpoint              string   `json:"authorization_endpoint"`
	TokenEndpoint                      string   `json:"token_endpoint"`
	UserinfoEndpoint                   string   `json:"userinfo_endpoint"`
	EndSessionEndpoint                 string   `json:"end_session_endpoint"`
	CheckSessionIframe                 string   `json:"check_session_iframe"`
	RevocationEndpoint                 string   `json:"revocation_endpoint"`
	IntrospectionEndpoint              string   `json:"introspection_endpoint"`
	DeviceAuthorizationEndpoint        string   `json:"device_authorization_endpoint"`
	FrontchannelLogoutSupported        bool     `json:"frontchannel_logout_supported"`
	FrontchannelLogoutSessionSupported bool     `json:"frontchannel_logout_session_supported"`
	BackchannelLogoutSupported         bool     `json:"backchannel_logout_supported"`
	BackchannelLogoutSessionSupported  bool     `json:"backchannel_logout_session_supported"`
	ScopesSupported                    []string `json:"scopes_supported"`
	ClaimsSupported                    []string `json:"claims_supported"`
	GrantTypesSupported                []string `json:"grant_types_supported"`
	ResponseTypesSupported             []string `json:"response_types_supported"`
	ResponseModesSupported             []string `json:"response_modes_supported"`
	TokenEndpointAuthMethodsSupported  []string `json:"token_endpoint_auth_methods_supported"`
	SubjectTypesSupported              []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported   []string `json:"id_token_signing_alg_values_supported"`
	CodeChallengeMethodsSupported      []string `json:"code_challenge_methods_supported"`
}

/*
{
	"issuer": "https://graphqlplay22.azurewebsites.net",
	"jwks_uri": "https://graphqlplay22.azurewebsites.net/.well-known/openid-configuration/jwks",
	"authorization_endpoint": "https://graphqlplay22.azurewebsites.net/connect/authorize",
	"token_endpoint": "https://graphqlplay22.azurewebsites.net/connect/token",
	"userinfo_endpoint": "https://graphqlplay22.azurewebsites.net/connect/userinfo",
	"end_session_endpoint": "https://graphqlplay22.azurewebsites.net/connect/endsession",
	"check_session_iframe": "https://graphqlplay22.azurewebsites.net/connect/checksession",
	"revocation_endpoint": "https://graphqlplay22.azurewebsites.net/connect/revocation",
	"introspection_endpoint": "https://graphqlplay22.azurewebsites.net/connect/introspect",
	"device_authorization_endpoint": "https://graphqlplay22.azurewebsites.net/connect/deviceauthorization",
	"frontchannel_logout_supported": true,
	"frontchannel_logout_session_supported": true,
	"backchannel_logout_supported": true,
	"backchannel_logout_session_supported": true,
	"scopes_supported": ["openid", "profile", "appIdentity", "offline_access"],
	"claims_supported": ["sub", "name", "family_name", "given_name", "middle_name", "nickname", "preferred_username", "profile", "picture", "website", "gender", "birthdate", "zoneinfo", "locale", "updated_at"],
	"grant_types_supported": ["authorization_code", "client_credentials", "refresh_token", "implicit", "urn:ietf:params:oauth:grant-type:device_code", "arbitrary_resource_owner", "arbitrary_identity", "arbitrary_no_subject"],
	"response_types_supported": ["code", "token", "id_token", "id_token token", "code id_token", "code token", "code id_token token"],
	"response_modes_supported": ["form_post", "query", "fragment"],
	"token_endpoint_auth_methods_supported": ["client_secret_basic", "client_secret_post"],
	"subject_types_supported": ["public"],
	"id_token_signing_alg_values_supported": ["RS256"],
	"code_challenge_methods_supported": ["plain", "S256"]
}
*/
