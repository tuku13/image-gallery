package pages

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
)

type RegisterPageData struct {
	Context *auth.JwtCustomClaims
}

func RegisterPage(c echo.Context) error {
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

	data := RegisterPageData{
		Context: context,
	}
	return c.Render(200, "register", data)
}
