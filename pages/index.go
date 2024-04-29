package pages

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/auth"
	"github.com/tuku13/image-gallery/constants"
	"github.com/tuku13/image-gallery/db/image"
	"github.com/tuku13/image-gallery/db/user"
)

type FormattedImage struct {
	Id           string
	Title        string
	Date         string
	Url          string
	UploaderName string
	UserId       string
}

type IndexPageData struct {
	Context *auth.JwtCustomClaims
	Images  []FormattedImage
	Order   string
}

func IndexPage(c echo.Context) error {
	var context *auth.JwtCustomClaims
	token, ok := c.Get(constants.CONTEXT_KEY).(*jwt.Token)
	if ok && token != nil {
		if claims, ok := token.Claims.(*auth.JwtCustomClaims); ok {
			context = claims
		}
	}

	var images []image.DbImage
	var err error
	query := c.QueryParam("query")
	orderBy := c.QueryParam("order_by")
	if orderBy == "title" {
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

	formattedImages := make([]FormattedImage, len(images))
	for i, img := range images {
		uploaderName := "Unknown"
		for _, dbUser := range users {
			if dbUser.Id == img.UserId {
				uploaderName = dbUser.Name
				break
			}
		}
		formattedImages[i] = FormattedImage{
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
