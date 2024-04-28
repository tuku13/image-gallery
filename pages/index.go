package pages

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/db/image"
)

type FormattedImage struct {
	ID           string
	Title        string
	Date         string
	Url          string
	UploaderName string
}

type IndexPageData struct {
	Context *auth.JwtCustomClaims
	Images  []FormattedImage
}

func IndexPage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
			context = claims
		}
	}

	images, err := image.GetImagesOrderByDate("")
	if err != nil {
		return c.String(500, "Internal Server Error")
	}

	formattedImages := make([]FormattedImage, len(images))
	for i, img := range images {
		formattedImages[i] = FormattedImage{
			ID:           img.Id,
			Title:        img.Title,
			Date:         img.UploadTime.Format("2006-01-02 15:04"),
			Url:          "/blob/" + img.BlobId,
			UploaderName: "TODO",
		}
	}

	pageData := IndexPageData{
		Context: context,
		Images:  formattedImages,
	}

	return c.Render(200, "index", pageData)
}
