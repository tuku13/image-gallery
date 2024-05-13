package api

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/database"
)

type Server struct {
	echo *echo.Echo
	db   database.Service
}

func NewServer() *Server {
	return &Server{
		echo: echo.New(),
		db:   database.New(),
	}
}

func (s *Server) Start() {
	s.echo.Renderer = newTemplate()
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	s.echo.Static("/static", "static")

	s.registerHandlers()

	s.echo.Logger.Fatal(s.echo.Start(":" + constants.PORT))
}

func (s *Server) registerHandlers() {
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(constants.JWT_SECRET),
		TokenLookup: "cookie:auth",
		ContextKey:  constants.CONTEXT_KEY,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
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

	public := s.echo.Group("")
	public.Use(checkJwtMiddleware)
	public.GET("/", s.handleIndexPage)
	public.GET("/login", s.handleLoginPage)
	public.POST("/auth/login", s.handleLogin)
	public.GET("/register", s.handleRegisterPage)
	public.POST("/auth/register", s.handleRegister)
	public.GET("/blob/:id", s.handleGetBlob)
	public.GET("/images/:id", s.handleImagePage)
	public.GET("/ordered_images", s.handleGetImages)

	private := s.echo.Group("")
	private.Use(jwtMiddleware)
	private.GET("/upload", s.handleUploadPage)
	private.POST("/upload-page/select-image", s.handleSelectImage)
	private.POST("/upload-page/deselect-image", s.handleDeselectImage)
	private.POST("/upload", s.handleUploadImage)
	private.POST("/auth/logout", s.handleLogout)
	private.DELETE("/images/:id", s.handleDeleteImage)
}
