package pages

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
)

type LoginPageData struct {
	Context *auth.JwtCustomClaims
}

func LoginPage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
			context = claims
		}
	}

	if context != nil {
		return c.Redirect(302, "/")
	}

	data := LoginPageData{
		Context: context,
	}
	return c.Render(200, "login", data)
}
