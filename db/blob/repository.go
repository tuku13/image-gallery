package blob

import (
	"github.com/jmoiron/sqlx"
	"github.com/tuku13/image-gallery/db"
	"log"
)

const schema = `
	CREATE TABLE IF NOT EXISTS blobs(
	    id UUID PRIMARY KEY,
	    data BYTEA
	);
`

type DbBlob struct {
	Id   string `db:"id"`
	Data []byte `db:"data"`
}

func InitDb() {
	db.Db.MustExec(schema)
}

func GetBlob(id string) (*DbBlob, error) {
	blob := DbBlob{}
	err := db.Db.Get(&blob, "SELECT * FROM blobs WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &blob, nil
}

func InsertBlobTx(tx *sqlx.Tx, blob *DbBlob) {
	tx.MustExec("INSERT INTO blobs (id, data) VALUES ($1, $2)", blob.Id, blob.Data)
}

func InsertBlob(blob *DbBlob) error {
	_, err := db.Db.Exec("INSERT INTO blobs (id, data) VALUES ($1, $2)", blob.Id, blob.Data)
	if err != nil {
		return err
	}
	return nil
}
