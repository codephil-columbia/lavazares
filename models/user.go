package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

//User metadata that is stored in the database
type User struct {
	UID       string     `json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
}

type UserSession struct {
	SessionID string
	UserID    string
	LoginTime time.Time
}

func NewUserSession(userID string) *UserSession {
	return &UserSession{
		xid.New().String(),
		userID,
		time.Now(),
	}
}

func NewUser(fields []byte) (string, error) {
	u := User{}
	err := json.Unmarshal(fields, &u)
	if err != nil {
		return "", err
	}

	u.UID = xid.New().String()
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return "", err
	}

	u.Password = hashedPassword
	result := db.QueryRowx("INSERT INTO users(UID, Username, Email, Password) VALUES($1, $2, $3, $4)",
		u.UID, u.Username, u.Email, u.Password).Err()
	if result != nil {
		return "", err
	}

	return u.UID, err
}

//AutheticateUser authenticates and returns a user
func AutheticateUser(req []byte) (*User, error) {
	u := User{}
	userAuthRequest := make(map[string]string)
	json.Unmarshal(req, &userAuthRequest)

	fmt.Println(userAuthRequest)

	err := db.QueryRowx("SELECT * FROM users WHERE email=$1",
		userAuthRequest["email"]).StructScan(&u)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(userAuthRequest["password"])); err != nil {
		return nil, err
	}

	return &u, nil
}

//HashPassword hashes a password and returns hashed password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Printf("error hashing password %s", err)
		return "", err
	}

	return string(bytes), nil
}
