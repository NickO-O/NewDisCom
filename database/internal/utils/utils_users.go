package utils

import (
	"database/sql"
	"db/internal/logger"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Login    string
	Password string
}

func NewUser(login, passoword string) (*User, error) {
	h, err := generate(passoword)
	if err != nil {
		return nil, err
	}

	return &User{Login: login, Password: h, Id: int(uuid.New().ID())}, nil
}

func generate(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func InsertUser(user *User) error {
	mu.Lock()
	defer mu.Unlock()
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		logger.Log.Error("Cannot open file data.sql")
		return err
	}

	_, err = db.Exec("INSERT INTO Users (id, login, Password) values ($1, $2, $3)", user.Id, user.Login, user.Password)
	if err != nil {
		fmt.Println(err, ", Insert user")
	}
	return nil
}

func SelectUserById(Id int) (User, error) {
	mu.Lock()
	defer mu.Unlock()
	db, err := sql.Open("sqlite3", dbname)

	if err != nil {
		logger.Log.Error("Cannot open file data.sql")
		return User{}, err
	}
	defer db.Close()
	row := db.QueryRow("select * from Users where id = $1", Id)
	usr, _ := NewUser("", "")
	err = row.Scan(&usr.Id, &usr.Login, &usr.Password)
	if err != nil {
		return User{}, err
	}
	return *usr, nil
}

func SelectUserByName(Name string) (User, error) {
	mu.Lock()
	defer mu.Unlock()
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		logger.Log.Error("Cannot open file data.sql")
		return User{}, err
	}
	defer db.Close()
	row := db.QueryRow("select * from Users where login = $1", Name)
	usr, _ := NewUser("", "")
	err = row.Scan(&usr.Id, &usr.Login, &usr.Password)

	if err != nil {
		return User{}, err
	}
	return *usr, nil
}

func Getallusers() ([]User, error) {
	mu.Lock()
	defer mu.Unlock()
	db, err := sql.Open("sqlite3", dbname)

	if err != nil {
		return nil, err
	}
	defer db.Close()
	all := make([]User, 0)
	rows, err := db.Query("SELECT * from Users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		usr, _ := NewUser("", "")
		err = rows.Scan(&usr.Id, &usr.Login, &usr.Password)
		if err != nil {
			logger.Log.Errorf("%s file: database.go func: GetAll", err.Error())
			return nil, err
		}

		all = append(all, *usr)
	}
	return all, nil

}
