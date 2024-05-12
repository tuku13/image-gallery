package database

import "time"

const userSchema = `
	CREATE TABLE IF NOT EXISTS users(
		id UUID PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT
	);
`
const blobSchema = `
	CREATE TABLE IF NOT EXISTS blobs(
	    id UUID PRIMARY KEY,
	    data BYTEA
	);
`
const imageSchema = `
	CREATE TABLE IF NOT EXISTS images(
	    id UUID PRIMARY KEY,
	    title varchar(40),
	    user_Id UUID REFERENCES users(id),
	    blob_id UUID REFERENCES blobs(id) ON DELETE CASCADE,
	    upload_time timestamp
	);
`

type Blob struct {
	Id   string `db:"id"`
	Data []byte `db:"data"`
}
type Image struct {
	Id         string    `db:"id"`
	Title      string    `db:"title"`
	UserId     string    `db:"user_id"`
	BlobId     string    `db:"blob_id"`
	UploadTime time.Time `db:"upload_time"`
}
type User struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
	Email    string `db:"email"`
}
