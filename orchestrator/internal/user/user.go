package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Login    string
	Password string
}

func NewUser(login, passoword string) (*User, error) {
	h, err := Generate(passoword)
	if err != nil {
		return nil, err
	}

	return &User{Login: login, Password: h, Id: int(uuid.New().ID())}, nil
}

func Generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}
