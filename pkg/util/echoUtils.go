package util

import (
	"fmt"

	echo "github.com/labstack/echo/v4"
)

func GetBaseUrl(c echo.Context) (baseUrl string) {
	request := c.Request()
	host := request.Host
	scheme := c.Scheme()
	baseUrl = fmt.Sprintf("%s://%s", scheme, host)
	return
}
