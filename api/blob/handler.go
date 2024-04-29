package api

import (
	"github.com/labstack/echo/v4"
	"github.com/tuku13/image-gallery/db/blob"
)

func GetBlob(c echo.Context) error {
	blobId := c.Param("id")

	blobData, err := blob.GetBlob(blobId)
	if err != nil {
		return c.String(404, "Image blob not found with id "+blobId)
	}

	c.Response().Header().Set("Cache-Control", "max-age=86400")
	return c.Blob(200, "image/webp", blobData.Data)
}
