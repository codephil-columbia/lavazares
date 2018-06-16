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
  FirstName  string     `json:"firstname"`
  LastName   string     `json:"lastname"`
	Username   string     `json:"username"`
  Firstname  string     `json:"firstname"` 
  Lastname   string     `json:"lastname"`
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
	res, err := db.Exec("INSERT INTO instructors(Gender, DOB, SchoolType, SchoolName, UID) VALUES($1, $2, $3, $4, $5)", s.Gender, s.DOB, s.SchoolType, s.SchoolName, u.UID)
	if res != nil {
		return err
	}

	return nil
}

func IsUsernameValid(req []byte) (bool, error) {
  var valid bool 
	body := make(map[string]string)
  err := json.Unmarshal(req, &body)
  if err != nil {
    return false, err
  }

  err = db.Get(&valid, "SELECT COUNT(*) FROM users WHERE username=$1 LIMIT 1", body["username"])
  if err != nil {
    return false, err
  }
  return !valid, err
}

func NewUser(fields []byte) (string, error) {

	u := User{}
	err := json.Unmarshal(fields, &u)
	if err != nil {
		return "", err
	}

	u.UID = xid.New().String()
	u.Password, err := hashPassword(u.Password)
	if err != nil {
		return "", err
	}

	result := db.QueryRowx("INSERT INTO users(UID, Firstname, Lastname, Username, Email, Password, Occupation) VALUES($1, $2, $3, $4, $5, $6, $7)",
		u.UID, u.FirstName, u.LastName, u.Username, u.Email, u.Password, u.Occupation).Err()
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

//AuthenticateUser authenticates and returns a user
func AuthenticateUser(userAuthRequest []byte) (*User, error) {
  u, u2 := User{}, User{}
  //userdata := make(map[string]interface{})
	err := json.Unmarshal(userAuthRequest, &u)
  if err != nil {
    return nil, err
  }

  err = db.QueryRowx("SELECT * FROM users WHERE username=$1", u.Username).StructScan(&u2)
  if err != nil {
    return nil, err
  }
  if err := bcrypt.CompareHashAndPassword([]byte(u2.Password), []byte(u.Password)); err != nil {
    return nil, err
  }

  return &u, nil
}

// EditPassword edits (rehashes) the password of an existing user
func EditPassword(req []byte) error {
  u := User{}
  body := make(map[string]string)
  json.Unmarshal(req, &body)

  fmt.Println(body)

  // TODO should grab username from session?
  err := db.QueryRowx("SELECT * FROM users WHERE username=$1 LIMIT 1", body["username"]).StructScan(&u)
  if err != nil {
    return err
  }

  hashedPassword, err := hashPassword(body["password"])
  if err != nil {
    return err
  }

  u.Password = hashedPassword
  result := db.QueryRowx("UPDATE users SET password=$1 WHERE username=$2", u.Password, u.Username)
  if result != nil {
    return err
  }
  
  return nil
}

// HashPassword hashes a password and returns hashed password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Printf("error hashing password %s", err)
		return "", err
	}

	return string(bytes), nil
}
