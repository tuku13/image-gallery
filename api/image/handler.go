package api

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/db/image"
)

func DeleteImage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
			context = claims
		}
	}
	if context == nil {
		return c.String(403, "Forbidden")
	}

	imageId := c.Param("id")
	dbImage, _ := image.GetImage(imageId)
	if dbImage == nil {
		return c.String(404, "Image not found")
	}

	if dbImage.UserId != context.UserId {
		return c.String(403, "Forbidden")
	}

	err := image.DeleteImage(imageId)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(200, "OK")
}
