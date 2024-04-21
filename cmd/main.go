package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e.Static("/static", "static")

	e.GET("/", pages.IndexPage)

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
