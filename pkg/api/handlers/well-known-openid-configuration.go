package handlers

import (
	"artificer/pkg/api/renderings"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func WellKnownOpenidConfiguration(c echo.Context) error {
	request := c.Request()
	host := request.Host
	scheme := c.Scheme()
	domain := fmt.Sprintf("%s://%s", scheme, host)
	issuer := domain

	jwks_uri := fmt.Sprintf("%s/.well-known/openid-configuration/jwks", domain)
	authorization_endpoint := fmt.Sprintf("%s/connect/authorize", domain)
	token_endpoint := fmt.Sprintf("%s/connect/token", domain)
	userinfo_endpoint := fmt.Sprintf("%s/connect/userinfo", domain)
	end_session_endpoint := fmt.Sprintf("%s/connect/endsession", domain)
	check_session_iframe := fmt.Sprintf("%s/connect/checksession", domain)
	revocation_endpoint := fmt.Sprintf("%s/connect/revocation", domain)
	introspection_endpoint := fmt.Sprintf("%s/connect/introspect", domain)
	device_authorization_endpoint := fmt.Sprintf("%s/connect/deviceauthorization", domain)

	resp := renderings.WellKnownOpenidConfigurationResponse{
		Issuer:                             issuer,
		JwksURI:                            jwks_uri,
		AuthorizationEndpoint:              authorization_endpoint,
		TokenEndpoint:                      token_endpoint,
		UserinfoEndpoint:                   userinfo_endpoint,
		EndSessionEndpoint:                 end_session_endpoint,
		CheckSessionIframe:                 check_session_iframe,
		RevocationEndpoint:                 revocation_endpoint,
		IntrospectionEndpoint:              introspection_endpoint,
		DeviceAuthorizationEndpoint:        device_authorization_endpoint,
		FrontchannelLogoutSupported:        false,
		FrontchannelLogoutSessionSupported: false,
		BackchannelLogoutSupported:         false,
		BackchannelLogoutSessionSupported:  false,
	}

	return c.JSON(http.StatusOK, resp)
}
