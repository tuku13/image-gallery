package api

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/database"
	"github.com/tuku13/image-gallery/types"
)

type ImagePageData struct {
	Context *JwtCustomClaims
	Image   types.FormattedImage
}
type UploadPageData struct {
	Context *JwtCustomClaims
}
type IndexPageData struct {
	Context *JwtCustomClaims
	Images  []types.FormattedImage
	Order   string
}
type LoginPageData struct {
	Context *JwtCustomClaims
}
type RegisterPageData struct {
	Context *JwtCustomClaims
}

func (s *Server) handleImagePage(c echo.Context) error {
	var context *JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*JwtCustomClaims); ok {
			context = claims
		}
	}

	imageId := c.Param("id")
	if imageId == "" {
		return c.String(404, "Not Found")
	}

	dbImage, err := s.db.GetImage(imageId)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	dbUser, err := s.db.GetUserById(dbImage.UserId)
	if err != nil {
		dbUser = &database.User{
			Id:       "Unknown",
			Name:     "Unknown",
			Email:    "Unknown",
			Password: "Unknown",
		}
	}

	pageData := ImagePageData{
		Context: context,
		Image: types.FormattedImage{
			Id:           dbImage.Id,
			Title:        dbImage.Title,
			Date:         dbImage.UploadTime.Format("2006-01-02 15:04:05"),
			Url:          "/blob/" + dbImage.BlobId,
			UploaderName: dbUser.Name,
			UserId:       dbImage.UserId,
		},
	}

	return c.Render(200, "image-page", pageData)
}
func (s *Server) handleUploadPage(c echo.Context) error {
	var context *JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*JwtCustomClaims); ok {
			context = claims
		}
	}

	data := UploadPageData{
		Context: context,
	}
	return c.Render(200, "upload", data)
}
func (s *Server) handleIndexPage(c echo.Context) error {
	var context *JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*JwtCustomClaims); ok {
			context = claims
		}
	}

	var images []database.Image
	var err error
	query := c.QueryParam("query")
	orderBy := c.QueryParam("order_by")
	if orderBy == "title" {
		images, err = s.db.GetImagesOrderByTitle(query)
	} else {
		images, err = s.db.GetImagesOrderByDate(query)
	}
	if err != nil {
		return c.String(500, "Failed to get images")
	}

	userIds := make(map[string]struct{})
	for _, img := range images {
		userIds[img.UserId] = struct{}{}
	}
	uniqueUserIds := make([]string, 0, len(userIds))
	for id := range userIds {
		uniqueUserIds = append(uniqueUserIds, id)
	}

	users := make([]database.User, len(uniqueUserIds))
	for i, id := range uniqueUserIds {
		dbUser, err := s.db.GetUserById(id)
		if err != nil {
			users[i] = database.User{
				Id:       id,
				Name:     "Unknown",
				Email:    "Unknown",
				Password: "Unknown",
			}
		} else {
			users[i] = *dbUser
		}
	}

	formattedImages := make([]types.FormattedImage, len(images))
	for i, img := range images {
		uploaderName := "Unknown"
		for _, dbUser := range users {
			if dbUser.Id == img.UserId {
				uploaderName = dbUser.Name
				break
			}
		}
		formattedImages[i] = types.FormattedImage{
			Id:           img.Id,
			Title:        img.Title,
			Date:         img.UploadTime.Format("2006-01-02 15:04"),
			Url:          "/blob/" + img.BlobId,
			UploaderName: uploaderName,
			UserId:       img.UserId,
		}
	}

	pageData := IndexPageData{
		Context: context,
		Images:  formattedImages,
		Order:   orderBy,
	}

	return c.Render(200, "index", pageData)
}
func (s *Server) handleLoginPage(c echo.Context) error {
	var context *JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*JwtCustomClaims); ok {
			context = claims
		}
	}

	if context != nil {
		return c.Redirect(302, "/")
	}

	data := LoginPageData{
		Context: context,
	}
	return c.Render(200, "login", data)
}
func (s *Server) handleRegisterPage(c echo.Context) error {
	var context *JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*JwtCustomClaims); ok {
			context = claims
		}
	}

	if context != nil {
		return c.Redirect(302, "/")
	}

	data := RegisterPageData{
		Context: context,
	}
	return c.Render(200, "register", data)
}
