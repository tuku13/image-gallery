package auth

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type User struct {
	Id    string
	Name  string
	Email string
}

type JWTSession struct {
	Id   string
	User User
}

type Context struct {
	echo.Context
	Session *JWTSession
}

func ContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//authCookie, _ := c.Cookie("auth")

		// TODO: Get session from auth cookie
		cc := &Context{
			Context: c,
		}

		return next(cc)
	}
}

func RequireContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println(c.Get("context"))
		// TODO check context is defined or not
		return next(c)
	}
}
