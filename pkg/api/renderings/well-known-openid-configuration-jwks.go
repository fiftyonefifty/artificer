package renderings

type JwkResponse struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	E   string `json:"e"`
	N   string `json:"n"`
	Alg string `json:"alg"`
}

type WellKnownOpenidConfigurationJwksResponse struct {
	Keys []JwkResponse `json:"keys"`
}

/*
{
	"keys": [{
		"kty": "RSA",
		"use": "sig",
		"kid": "https://p7keyvalut.vault.azure.net/keys/P7IdentityServer4SelfSigned/8e87bfff06004b81a10e5acfbf5f852a",
		"e": "AAEAAQ",
		"n": "q9GgMKX_D-v4MTqmEyeMvTHRV_8d1SHdisL10XXCkbld980AkXU8s8_wR6_Skvjwx5HTRl2sL82qAK93uCJ0z1wgLYMl_0YU3dPZAw7UzJmpeR27f2beFy2GvFyzD8DtNXHv2pL5Rwh040LCqJMz-PNShxat3UK6H-PWa82uACxz9OWcP5nld0fLGqcah5yiNVFPuAxyyqIrsD22jj4f0kMd1YdIYDQcQF4r-3coHnxterXw461Qo2ruVRavPxkS7G9dELMw6UQK3AA9hpqgJ489xvk3Xe6O8aO4-Jy78LrGX2B0zlEnhkZJ2K30hGlJzlxbRWVeLAAsoxmNtcwTCQ",
		"alg": "RS256"
	}, {
		"kty": "RSA",
		"use": "sig",
		"kid": "https://p7keyvalut.vault.azure.net/keys/P7IdentityServer4SelfSigned/f7d87a46708c4f3d8dfe611e973745c3",
		"e": "AAEAAQ",
		"n": "41KzG6IBcFXCoMLOjcihqoKkLqL8M_cRI0k2HCTzoYZOpwdk8l1jIqSIX8BtczuCGEvIj17l9ykkoKmC1anWnNJxvIam_2dnciZNvIIqgivf0E-S-1J_Ushxz-Fw0okUknFzX1gxyWEAMWxNtscwiB3_Z3-SyiVZWMrTRCxnwa37E7QjkLFNHuxXd0q7Hy9h6EQ8cx1lPsDxItkOQb8wSWduxVot94gSry5aEe4eG6gKXk0b5fJSLF0OsDErdXcmvYw_8_sCReXgR3UfmXwkr4NVQOKP6XLFgpdzcML8fY5PqK0uJY_iFDYuR-jIbMsD-DOy4N5jSBUz-q8vdWjUww",
		"alg": "RS256"
	}]
}
*/
