package api

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/database"
	"github.com/tuku13/image-gallery/types"
	"io"
	"mime/multipart"
	"time"
)

func (s *Server) handleGetImages(c echo.Context) error {
	var images []database.Image
	var err error
	query := c.QueryParam("query")
	fmt.Println("OrderedImages", query)
	orderBy := c.QueryParam("order_by")
	if orderBy == "name" {
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

	return c.Render(200, "images", types.OrderedImages{Images: formattedImages})
}
func (s *Server) handleDeleteImage(c echo.Context) error {
	var context *JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*JwtCustomClaims); ok {
			context = claims
		}
	}
	if context == nil {
		return c.String(403, "Forbidden")
	}

	imageId := c.Param("id")
	dbImage, _ := s.db.GetImage(imageId)
	if dbImage == nil {
		return c.String(404, "Image not found")
	}

	if dbImage.UserId != context.UserId {
		return c.String(403, "Forbidden")
	}

	err := s.db.DeleteImage(imageId)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(200, "OK")
}

func (s *Server) handleGetBlob(c echo.Context) error {
	blobId := c.Param("id")

	blobData, err := s.db.GetBlob(blobId)
	if err != nil {
		return c.String(404, "Image blob not found with id "+blobId)
	}

	c.Response().Header().Set("Cache-Control", "max-age=86400")
	return c.Blob(200, "image/webp", blobData.Data)
}

func (s *Server) handleSelectImage(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	return c.Render(200, "image-uploader-control-filled", map[string]interface{}{
		"Filename": file.Filename,
	})
}
func (s *Server) handleDeselectImage(c echo.Context) error {
	return c.Render(200, "image-uploader", nil)
}
func (s *Server) handleUploadImage(c echo.Context) error {
	var context *JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*JwtCustomClaims); ok {
			context = claims
		}
	}
	if context == nil {
		return c.Redirect(302, "/login")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(400, "Title is required")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			c.Logger().Fatal(err)
		}
	}(src)

	data, err := io.ReadAll(src)
	if err != nil {
		return c.String(500, "Failed to read image")
	}

	blobId := uuid.New().String()
	dbBlob := &database.Blob{
		Id:   blobId,
		Data: data,
	}
	if err = s.db.InsertBlob(dbBlob); err != nil {
		return err
	}

	imageId := uuid.New().String()
	dbImage := &database.Image{
		Id:         imageId,
		Title:      title,
		BlobId:     blobId,
		UserId:     context.UserId,
		UploadTime: time.Now(),
	}
	if err = s.db.InsertImage(dbImage); err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/images/"+imageId)
	return c.String(200, "OK")
}
