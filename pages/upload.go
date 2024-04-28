package pages

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
)

type UploadPageData struct {
	Context *auth.JwtCustomClaims
}

func UploadPage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
			context = claims
		}
	}

	data := UploadPageData{
		Context: context,
	}
	return c.Render(200, "upload", data)
}

func SelectImagePost(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	return c.Render(200, "image-uploader-control-filled", map[string]interface{}{
		"Filename": file.Filename,
	})
}

func DeselectImagePost(c echo.Context) error {
	return c.Render(200, "image-uploader", nil)
}

func UploadImage(c echo.Context) error {
	title := c.FormValue("title")
	if title == "" {
		return c.String(400, "Title is required")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// TODO save the file to db and redirect to /images/:id

	return c.NoContent(200)
}
