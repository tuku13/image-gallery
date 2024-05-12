package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Service interface {
	GetBlob(id string) (*Blob, error)
	InsertBlob(blob *Blob) error

	GetImage(id string) (*Image, error)
	DeleteImage(id string) error
	InsertImage(image *Image) error
	GetImagesOrderByTitle(query string) ([]Image, error)
	GetImagesOrderByDate(query string) ([]Image, error)

	GetUserByName(name string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	InsertUser(user *User) error
}

type service struct {
	db *sqlx.DB
}

var dbInstance *service

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "image-gallery-db"
	}
	usr := os.Getenv("DB_USER")
	if usr == "" {
		usr = "test"
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = "test"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	uri := "user=" + usr + " dbname=" + dbname + " password=" + pass + " host=" + host + " port=" + port + " sslmode=disable"
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		panic(err)
	}

	dbInstance = &service{
		db: db,
	}
	dbInstance.mustInit()
	return &service{
		db: db,
	}
}
func (s *service) mustInit() {
	s.db.MustExec(userSchema)
	s.db.MustExec(blobSchema)
	s.db.MustExec(imageSchema)
}

func (s *service) GetBlob(id string) (*Blob, error) {
	blob := Blob{}
	err := s.db.Get(&blob, "SELECT * FROM blobs WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &blob, nil
}
func (s *service) InsertBlob(blob *Blob) error {
	_, err := s.db.Exec("INSERT INTO blobs (id, data) VALUES ($1, $2)", blob.Id, blob.Data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetImage(id string) (*Image, error) {
	image := Image{}
	err := s.db.Get(&image, "SELECT * FROM images WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &image, nil
}
func (s *service) DeleteImage(id string) error {
	_, err := s.db.Exec("DELETE FROM images WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) InsertImage(image *Image) error {
	_, err := s.db.Exec("INSERT INTO images (id, title, user_id, blob_id, upload_time) VALUES ($1, $2, $3, $4, $5)", image.Id, image.Title, image.UserId, image.BlobId, image.UploadTime)
	return err
}
func (s *service) GetImagesOrderByTitle(query string) ([]Image, error) {
	var images []Image
	var err error
	if query == "" {
		err = s.db.Select(&images, "SELECT * FROM images ORDER BY title DESC")
	} else {
		err = s.db.Select(&images, "SELECT * FROM images WHERE title ilike $1 ORDER BY title DESC", "%"+query+"%")
	}
	if err != nil {
		return nil, err
	}
	return images, nil
}
func (s *service) GetImagesOrderByDate(query string) ([]Image, error) {
	var images []Image
	var err error
	if query == "" {
		err = s.db.Select(&images, "SELECT * FROM images ORDER BY upload_time DESC")
	} else {
		err = s.db.Select(&images, "SELECT * FROM images WHERE title ilike $1 ORDER BY upload_time DESC", "%"+query+"%")
	}
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (s *service) GetUserByName(name string) (*User, error) {
	user := User{}
	err := s.db.Get(&user, "SELECT * FROM users WHERE name=$1", name)
	return &user, err
}
func (s *service) GetUserByEmail(email string) (*User, error) {
	user := User{}
	err := s.db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	return &user, err
}
func (s *service) GetUserById(id string) (*User, error) {
	user := User{}
	err := s.db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return &user, err
}
func (s *service) InsertUser(user *User) error {
	_, err := s.db.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)", user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
