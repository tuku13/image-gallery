package db

type User struct {
	Id    string
	Name  string
	Email string
}

func GetUser(email string) *User {
	return &User{
		Id:    "1",
		Name:  "John Doe",
		Email: email,
	}
}
