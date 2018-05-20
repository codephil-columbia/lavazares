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
	UID        string     `json:"-"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `json:"-"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Occupation string     `json:"occupation"`
}

type Student struct {
	Gender             string `json:"gender" db:"gender"`
	DOB                string `json:"dob" db:"dob"`
	SchoolYear         int    `json:"schoolyear" db:"schoolyear"`
	CurrentLessonID    string `json:"currentlessonid" db:"currentlessonid"`
	CurrentChapterName string `db:"currentchaptername"`
	UID                string `db:"uid"`
}

type Instructor struct {
	Gender     string `json:"gender" db:"gender"`
	DOB        string `json:"dob" db:"dob"`
	SchoolType string `json:"schooltype" db:"schooltype"`
	SchoolName string `json:"schoolname" db:"schoolname"`
	UID        string `json:"uid" db:"uid"`
}

func NewStudent(fields []byte, u User) error {
	s := Student{}
	err := json.Unmarshal(fields, &s)
	if err != nil {
		return err
	}
	fmt.Println(s)
	res, err := db.Exec("INSERT INTO students(Gender, DOB, SchoolYear, CurrentLessonID, CurrentChapterName UID) VALUES($1, $2, $3, $4, $5, $6)",
		s.Gender, s.DOB, s.SchoolYear, s.CurrentLessonID, "Chapter 0: The Basics", u.UID)
	if res != nil {
		return err
	}

	return nil
}

func NewInstructor(fields []byte, u User) error {
	s := Instructor{}
	err := json.Unmarshal(fields, &s)
	if err != nil {
		return err
	}
	fmt.Println(s)
	res, err := db.Exec("INSERT INTO instructors(Gender, DOB, SchoolType, SchoolName, UID) VALUES($1, $2, $3, $4, $5)",
		s.Gender, s.DOB, s.SchoolType, s.SchoolName, u.UID)
	if res != nil {
		return err
	}

	return nil
}

func UsernameExists(username string) bool {
  result := db.QueryRowx("SELECT COUNT(*) FROM users WHERE username='$1' LIMIT 1", username);
  return result
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

	fmt.Println(u)

	u.Password = hashedPassword
	result := db.QueryRowx("INSERT INTO users(UID, Username, Email, Password, Occupation) VALUES($1, $2, $3, $4, $5)",
		u.UID, u.Username, u.Email, u.Password, u.Occupation).Err()
	if result != nil {
		return "", err
	}

	switch u.Occupation {
	case "student":
		NewStudent(fields, u)
	case "instructor":
		NewInstructor(fields, u)
	}

	return u.UID, err
}

func UpdateModel(modelName, field, value, identifier, identifierVal string) error {
	stmt := fmt.Sprintf("UPDATE %s SET %s='%s' WHERE %s='%s'", modelName, field, value, identifier, identifierVal)
	fmt.Println(stmt)
	result, err := db.Exec(stmt)
	if result != nil {
		return err
	}

	fmt.Println(result)

	return nil
}

func GetStudent(uid string) (*Student, error) {
	s := Student{}
	err := db.QueryRowx("SELECT * FROM students WHERE uid=$1", uid).StructScan(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
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
