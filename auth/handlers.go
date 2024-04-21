package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/db"
	"net/http"
	"time"
)

type LoginRequestForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type RegisterRequestForm struct {
	Username      string `form:"username"`
	Email         string `form:"email"`
	Password      string `form:"password"`
	PasswordAgain string `form:"password_again"`
}

func createJWT(user *db.User, expires time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": expires.Unix(),
	})
	tokenString, err := token.SignedString([]byte(constants.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginPost(c echo.Context) error {
	var loginRequest LoginRequestForm
	if err := c.Bind(&loginRequest); err != nil {
		fmt.Println(err)
		return c.String(400, "Bad Request")
	}

	user := db.GetUser(loginRequest.Email)
	if user == nil {
		return c.String(401, "Unauthorized")
	}

	expires := time.Now().Add(1 * time.Hour)
	jwtString, err := createJWT(user, expires)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not create JWT token")
	}

	// Set the token in a cookie
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = jwtString
	cookie.Expires = expires
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.NoContent(200)
}

func RegisterPost(c echo.Context) error {
	var registerRequest RegisterRequestForm
	if err := c.Bind(&registerRequest); err != nil {
		return c.String(400, "Bad Request")
	}

	if err := c.Validate(registerRequest); err != nil {
		return c.String(400, "Bad Request")
	}

	return c.String(200, "Login Request Post")
}

func LogoutPost(c echo.Context) error {
	var cookie *http.Cookie
	cookie, err := c.Cookie("auth")
	if err != nil {
		cookie = &http.Cookie{
			Name: "auth",
		}
	}

	cookie.MaxAge = -1
	cookie.Value = ""
	c.SetCookie(cookie)

	return c.NoContent(200)
}
