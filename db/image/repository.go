package image

import (
	"github.com/jmoiron/sqlx"
	"github.com/tuku13/image-gallery/db"
	"time"
)

const schema = `
	CREATE TABLE IF NOT EXISTS images(
	    id UUID PRIMARY KEY,
	    title varchar(40),
	    user_Id UUID REFERENCES users(id),
	    blob_id UUID REFERENCES blobs(id) ON DELETE CASCADE,
	    upload_time timestamp
	);
`

type DbImage struct {
	Id         string    `db:"id"`
	Title      string    `db:"title"`
	UserId     string    `db:"user_id"`
	BlobId     string    `db:"blob_id"`
	UploadTime time.Time `db:"upload_time"`
}

func InitDb() {
	db.Db.MustExec(schema)
}

func InsertImageTx(tx *sqlx.Tx, image *DbImage) {
	tx.MustExec("INSERT INTO images (id, title, user_id, blob_id, upload_time) VALUES ($1, $2, $3, $4, $5)", image.Id, image.Title, image.UserId, image.BlobId, image.UploadTime)
}

func GetImage(id string) (*DbImage, error) {
	image := DbImage{}
	err := db.Db.Get(&image, "SELECT * FROM images WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func DeleteImage(id string) error {
	_, err := db.Db.Exec("DELETE FROM images WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func GetImagesOrderByTitle(query string) ([]DbImage, error) {
	var images []DbImage
	var err error
	if query == "" {
		err = db.Db.Select(&images, "SELECT * FROM images ORDER BY title DESC")
	} else {
		err = db.Db.Select(&images, "SELECT * FROM images WHERE title ilike $1 ORDER BY title DESC", "%"+query+"%")
	}
	if err != nil {
		return nil, err
	}
	return images, nil
}

func GetImagesOrderByDate(query string) ([]DbImage, error) {
	var images []DbImage
	var err error
	if query == "" {
		err = db.Db.Select(&images, "SELECT * FROM images ORDER BY upload_time DESC")
	} else {
		err = db.Db.Select(&images, "SELECT * FROM images WHERE title ilike $1 ORDER BY upload_time DESC", "%"+query+"%")
	}
	if err != nil {
		return nil, err
	}
	return images, nil
}
