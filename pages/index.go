package pages

import "github.com/labstack/echo/v4"

type FormattedImage struct {
	ID           string
	Title        string
	Date         string
	Url          string
	UploaderName string
}

type data struct {
	Images []FormattedImage
}

func IndexPage(c echo.Context) error {

	pageData := data{
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
