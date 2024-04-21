package pages

import "github.com/labstack/echo/v4"

func UploadPage(c echo.Context) error {
	return c.Render(200, "upload", nil)
}
