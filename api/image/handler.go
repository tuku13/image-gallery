package api

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/db/image"
	"github.com/tuku13/image-gallery/db/user"
	"github.com/tuku13/image-gallery/pages"
)

type orderedImagesType struct {
	Images []pages.FormattedImage
}

func OrderedImages(c echo.Context) error {
	var images []image.DbImage
	var err error
	query := c.QueryParam("query")
	fmt.Println("OrderedImages", query)
	orderBy := c.QueryParam("order_by")
	if orderBy == "name" {
		images, err = image.GetImagesOrderByTitle(query)
	} else {
		images, err = image.GetImagesOrderByDate(query)
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

	users := make([]user.DbUser, len(uniqueUserIds))
	for i, id := range uniqueUserIds {
		dbUser, err := user.GetUserById(id)
		if err != nil {
			users[i] = user.DbUser{
				Id:       id,
				Name:     "Unknown",
				Email:    "Unknown",
				Password: "Unknown",
			}
		} else {
			users[i] = *dbUser
		}
	}

	formattedImages := make([]pages.FormattedImage, len(images))
	for i, img := range images {
		uploaderName := "Unknown"
		for _, dbUser := range users {
			if dbUser.Id == img.UserId {
				uploaderName = dbUser.Name
				break
			}
		}
		formattedImages[i] = pages.FormattedImage{
			Id:           img.Id,
			Title:        img.Title,
			Date:         img.UploadTime.Format("2006-01-02 15:04"),
			Url:          "/blob/" + img.BlobId,
			UploaderName: uploaderName,
			UserId:       img.UserId,
		}
	}

	return c.Render(200, "images", orderedImagesType{Images: formattedImages})
}

func DeleteImage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
			context = claims
		}
	}
	if context == nil {
		return c.String(403, "Forbidden")
	}

	imageId := c.Param("id")
	dbImage, _ := image.GetImage(imageId)
	if dbImage == nil {
		return c.String(404, "Image not found")
	}

	if dbImage.UserId != context.UserId {
		return c.String(403, "Forbidden")
	}

	err := image.DeleteImage(imageId)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(200, "OK")
}
