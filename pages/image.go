package pages

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/db/image"
	"github.com/tuku13/image-gallery/db/user"
)

type ImagePageData struct {
	Context *auth.JwtCustomClaims
	Image   FormattedImage
}

func ImagePage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
			context = claims
		}
	}

	imageId := c.Param("id")
	if imageId == "" {
		return c.String(404, "Not Found")
	}

	dbImage, err := image.GetImage(imageId)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	dbUser, err := user.GetUserById(dbImage.UserId)
	if err != nil {
		dbUser = &user.DbUser{
			Id:       "Unknown",
			Name:     "Unknown",
			Email:    "Unknown",
			Password: "Unknown",
		}
	}

	pageData := ImagePageData{
		Context: context,
		Image: FormattedImage{
			Id:           dbImage.Id,
			Title:        dbImage.Title,
			Date:         dbImage.UploadTime.Format("2006-01-02 15:04:05"),
			Url:          "/blob/" + dbImage.BlobId,
			UploaderName: dbUser.Name,
			UserId:       dbImage.UserId,
		},
	}

	return c.Render(200, "image-page", pageData)
}
