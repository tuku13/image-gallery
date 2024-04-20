package pages

import "github.com/labstack/echo/v4"

func IndexPage(c echo.Context) error {
	return c.Render(200, "index", nil)
}
