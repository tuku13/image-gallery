package pages

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
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

	pageData := IndexPageData{
		Context: context,
		Images: []FormattedImage{
			{ID: "1", Title: "FormattedImage 1FormattedImage 1FormattedImage 1v", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
			{ID: "1", Title: "FormattedImage 1", Date: "2024-04-21 13:50", UploaderName: "user123", Url: "https://placehold.co/600x400@2x.png"},
		},
	}

	return c.Render(200, "index", pageData)
}
