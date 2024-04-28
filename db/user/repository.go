package user

import (
	"github.com/tuku13/image-gallery/db"
)

const schema = `
	CREATE TABLE IF NOT EXISTS users(
		id UUID PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT
	);
`

type DbUser struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

func InitDb() {
	db.Db.MustExec(schema)
}

func GetUserByName(name string) (*DbUser, error) {
	user := DbUser{}
	err := db.Db.Get(&user, "SELECT * FROM users WHERE name=$1", name)
	return &user, err
}

func GetUserByEmail(email string) (*DbUser, error) {
	user := DbUser{}
	err := db.Db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	return &user, err
}

func GetUserById(id string) (*DbUser, error) {
	user := DbUser{}
	err := db.Db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return &user, err
}

func InsertUser(user *DbUser) error {
	_, err := db.Db.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)", user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
