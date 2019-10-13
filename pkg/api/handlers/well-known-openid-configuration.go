package handlers

import (
	"artificer/pkg/api/renderings"
	"artificer/pkg/util"
	"fmt"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	gocache "github.com/pmylund/go-cache"
)

var (
	cache    = gocache.New(24*time.Hour, time.Hour)
	cacheKey = "64e3422b-c321-4c4a-b03e-66c92f6f53b5"
)

func buildCachedResponse(c echo.Context) {
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

	cacheItem := renderings.WellKnownOpenidConfigurationResponse{
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
	cache.Set(cacheKey, cacheItem, gocache.NoExpiration)
}

// HealthCheck - Healthcheck Handler
func WellKnownOpenidConfiguration(c echo.Context) error {

	var cachedItem interface{}
	var found bool

	cachedItem, found = cache.Get(cacheKey)
	if !found {
		buildCachedResponse(c)
		cachedItem, found = cache.Get(cacheKey)
	}
	resp := cachedItem.(renderings.WellKnownOpenidConfigurationResponse)

	return c.JSON(http.StatusOK, resp)
}
