package entity

type User struct {
	Id   int
	Name string
	Key  []byte
}

func NewUser(name string, key []byte) User {
	return User{Name: name, Key: key}
}
