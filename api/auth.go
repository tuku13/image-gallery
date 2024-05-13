package api

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/database"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginRequestForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type RegisterRequestForm struct {
	Name          string `form:"name"`
	Email         string `form:"email"`
	Password      string `form:"password"`
	PasswordAgain string `form:"passwordAgain"`
}

type JwtCustomClaims struct {
	Name   string `json:"name"`
	UserId string `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *Server) handleLogin(c echo.Context) error {
	var loginRequest LoginRequestForm
	if err := c.Bind(&loginRequest); err != nil {
		fmt.Println(err)
		return c.String(400, "Bad Request")
	}

	loggedInUser, _ := s.db.GetUserByEmail(loginRequest.Email)
	fmt.Println(loggedInUser)
	if loggedInUser == nil {
		return c.String(401, "Unauthorized")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(loggedInUser.Password), []byte(loginRequest.Password)); err != nil {
		return c.String(401, "Unauthorized")
	}

	expires := time.Now().Add(1 * time.Hour)
	jwtString, err := createJWT(loggedInUser, expires)
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

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(200)
}
func (s *Server) handleRegister(c echo.Context) error {
	var registerRequest RegisterRequestForm
	if err := c.Bind(&registerRequest); err != nil {
		return c.String(400, "Bad Request")
	}

	if registerRequest.Name == "" {
		return c.String(400, "Username is required")
	}
	if registerRequest.Email == "" {
		return c.String(400, "Email is required")
	}
	if registerRequest.Password == "" {
		return c.String(400, "Password is required")
	}
	if registerRequest.Password != registerRequest.PasswordAgain {
		return c.String(400, "Passwords do not match")
	}

	registeredByMail, _ := s.db.GetUserByEmail(registerRequest.Email)
	if registeredByMail == nil {
		return c.String(400, "Email already registered")
	}
	registeredByUsername, _ := s.db.GetUserByName(registerRequest.Name)
	if registeredByUsername == nil {
		return c.String(400, "Username already registered")
	}

	hashedPassword, err := hashPassword(registerRequest.Password)
	if err != nil {
		return c.String(500, err.Error())
	}

	dbUser := database.User{
		Id:       uuid.New().String(),
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: hashedPassword,
	}
	err = s.db.InsertUser(&dbUser)
	if err != nil {
		return c.String(500, "Could not register user")
	}

	expires := time.Now().Add(1 * time.Hour)
	jwtString, err := createJWT(&dbUser, expires)
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

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(200)
}
func (s *Server) handleLogout(c echo.Context) error {
	var cookie *http.Cookie
	cookie, err := c.Cookie("auth")
	if err != nil {
		cookie = &http.Cookie{
			Name: "auth",
		}
	}

	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Value = ""
	cookie.Path = "/"
	c.SetCookie(cookie)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(200)
}

func createJWT(user *database.User, expires time.Time) (string, error) {
	claims := &JwtCustomClaims{
		Name:   user.Name,
		UserId: user.Id,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(constants.JWT_SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
