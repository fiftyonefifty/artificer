package handlers

import (
	"artificer/pkg/api/renderings"
	"artificer/pkg/util"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func WellKnownOpenidConfiguration(c echo.Context) error {

	issuer := util.GetBaseUrl(c)

	jwks_uri := fmt.Sprintf("%s/.well-known/openid-configuration/jwks", issuer)
	authorization_endpoint := fmt.Sprintf("%s/connect/authorize", issuer)
	token_endpoint := fmt.Sprintf("%s/connect/token", issuer)
	userinfo_endpoint := fmt.Sprintf("%s/connect/userinfo", issuer)
	end_session_endpoint := fmt.Sprintf("%s/connect/endsession", issuer)
	check_session_iframe := fmt.Sprintf("%s/connect/checksession", issuer)
	revocation_endpoint := fmt.Sprintf("%s/connect/revocation", issuer)
	introspection_endpoint := fmt.Sprintf("%s/connect/introspect", issuer)
	device_authorization_endpoint := fmt.Sprintf("%s/connect/deviceauthorization", issuer)

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
