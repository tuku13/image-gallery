package main

import (
	"fmt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/pages"
	"io"
)

import (
	"html/template"
)

func main() {
	e := echo.New()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(auth.ContextMiddleware)

	e.Static("/static", "static")

	public := e.Group("")
	public.GET("/", pages.IndexPage)
	public.POST("/auth/login", auth.LoginPost)
	public.POST("/auth/register", auth.RegisterPost)

	private := e.Group("")
	private.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(constants.JWT_SECRET),
		TokenLookup: "cookie:auth",
		ContextKey:  "context",
		ErrorHandler: func(c echo.Context, err error) error {
			// print cookies
			cookies := c.Cookies()
			for _, cookie := range cookies {
				fmt.Println(cookie.Name + ":" + cookie.Value)
			}
			return c.String(401, "Unauthorized")
		},
	}))

	private.GET("/upload", pages.UploadPage)
	private.POST("/auth/logout", auth.LogoutPost)

	e.Logger.Fatal(e.Start(":" + constants.PORT))
}

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/**/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}
