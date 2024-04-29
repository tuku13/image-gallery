package pages

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/db"
	"github.com/tuku13/image-gallery/db/blob"
	"github.com/tuku13/image-gallery/db/image"
	"io"
	"time"
)

type UploadPageData struct {
	Context *auth.JwtCustomClaims
}

func UploadPage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
			context = claims
		}
	}

	data := UploadPageData{
		Context: context,
	}
	return c.Render(200, "upload", data)
}

func SelectImagePost(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}

	return c.Render(200, "image-uploader-control-filled", map[string]interface{}{
		"Filename": file.Filename,
	})
}

func DeselectImagePost(c echo.Context) error {
	return c.Render(200, "image-uploader", nil)
}

func UploadImage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
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
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return c.String(500, "Failed to read image")
	}

	tx := db.Db.MustBegin()

	blobId := uuid.New().String()
	dbBlob := &blob.DbBlob{
		Id:   blobId,
		Data: data,
	}
	blob.InsertBlobTx(tx, dbBlob)

	imageId := uuid.New().String()
	dbImage := &image.DbImage{
		Id:         imageId,
		Title:      title,
		BlobId:     blobId,
		UserId:     context.UserId,
		UploadTime: time.Now(),
	}
	image.InsertImageTx(tx, dbImage)

	err = tx.Commit()
	if err != nil {
		return c.String(500, "Failed to save image")
	}

	c.Response().Header().Set("HX-Redirect", "/images/"+imageId)
	return c.String(200, "OK")
}
