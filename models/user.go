package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

//User metadata that is stored in the database
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UID      string `json:"uid"`
}

type UserSession struct {
	sessions.Session
	SessionID    string
	UserID       string
	LoginTime    time.Time
	LastSeenTime time.Time
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// func NewUserSession(userID string) (*UserSession, error) {
// 	session := UserSession{}
// 	session.UserID = userID
// 	session.SessionID = randomKey()
// 	session.LoginTime = time.Now()
// 	session.LastSeenTime = time.Now()
// }

//RetrieveUser retrieves a user given an email
func RetrieveUser(email string) (*User, error) {
	user := User{}
	err := db.QueryRow("select username, password, email, uid from users where email=$1", email).Scan(&user.Username, &user.Password, &user.Email, &user.UID)
	if err != nil {
		log.Printf("error scanning user: %s", err)
		return nil, err
	}
	return &user, nil
}

//InsertUser inserts a new user into the database
func InsertUser(user *User) error {
	if db == nil {
		fmt.Println("database is nil")
	}
	_, err := db.Query(
		"insert into users(username, password, email, uid) values($1, $2, $3, $4)",
		&user.Username,
		&user.Password,
		&user.Email,
		&user.UID,
	)

	if err != nil {
		log.Printf("error adding user to database: %s", err)
	}

	return nil
}

//AutheticateUser authenticates and returns a user
func AutheticateUser(req *LoginRequest) (*User, error) {
	user, err := RetrieveUser(req.Email)
	if err != nil {
		log.Printf("error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("passwords do not match: %s", err)
		return nil, err
	}

	return user, nil
}

//HashPassword hashes a password and returns hashed password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Printf("error hashing password %s", err)
		return "", err
	}

	return string(bytes), nil
}

func RandomKey() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
