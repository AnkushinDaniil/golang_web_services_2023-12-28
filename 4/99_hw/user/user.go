package user

type User struct {
	Id     int
	Name   string
	Age    int
	About  string
	Gender string
}

//easyjson:json
type Users []User
