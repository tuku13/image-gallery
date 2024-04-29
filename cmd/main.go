package main

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	api "github.com/tuku13/image-gallery/api/blob"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/db/blob"
	"github.com/tuku13/image-gallery/db/image"
	"github.com/tuku13/image-gallery/db/user"
	"github.com/tuku13/image-gallery/pages"
	"io"
)

import (
	"html/template"
)

func main() {
	user.InitDb()
	blob.InitDb()
	image.InitDb()

	e := echo.New()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/static", "static")

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(constants.JWT_SECRET),
		TokenLookup: "cookie:auth",
		ContextKey:  constants.CONTEXT_KEY,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
	})
	checkJwtMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("auth")
			if err == nil && cookie.Value != "" {
				return jwtMiddleware(next)(c)
			}
			return next(c)
		}
	}

	public := e.Group("")
	public.Use(checkJwtMiddleware)
	public.GET("/", pages.IndexPage)
	public.GET("/login", pages.LoginPage)
	public.POST("/auth/login", auth.LoginPost)
	public.GET("/register", pages.RegisterPage)
	public.POST("/auth/register", auth.RegisterPost)
	public.GET("/blob/:id", api.GetBlob)
	public.GET("/images/:id", pages.ImagePage)

	private := e.Group("")
	private.Use(jwtMiddleware)
	private.GET("/upload", pages.UploadPage)
	private.POST("/upload-page/select-image", pages.SelectImagePost)
	private.POST("/upload-page/deselect-image", pages.DeselectImagePost)
	private.POST("/upload", pages.UploadImage)
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
