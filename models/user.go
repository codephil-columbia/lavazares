package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

var (
	CURRENT_LESSON_ID = "d3f9c2a3-1edf-42a6-a24d-3a4ad4683036"
	CURRENT_CHAPTER_NAME = "Chapter 0: The Basics"
)

//User metadata that is stored in the database
type User struct {
	UID             string     `json:"-"`
	CreatedAt       time.Time  `json:"-"`
	UpdatedAt       time.Time  `json:"-"`
	DeletedAt       *time.Time `json:"-"`
	FirstName       string     `json:"firstname"`
	LastName        string     `json:"lastname"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	Occupation      string     `json:"occupation"`
	WhichOccupation string     `json:"whichOccupation"`
}

type Student struct {
	Gender             string `json:"gender" db:"gender"`
	DOB                string `json:"dob" db:"dob"`
	SchoolYear         string `json:"schoolyear" db:"schoolyear"`
	CurrentLessonID    string `json:"currentlessonid" db:"currentlessonid"`
	CurrentChapterName string `db:"currentchaptername"`
	UID                string `db:"uid"`
}

type Employed struct {
	Gender             string `json:"gender" db:"gender"`
	DOB                string `json:"dob" db:"dob"`
	CurrentLessonID    string `json:"currentlessonid" db:"currentlessonid"`
	CurrentChapterName string `db:"currentchaptername"`
	UID                string `db:"uid"`
	occupation string `json:"occupation" db:"occupation"`
}

type Instructor struct {
	Gender     string `json:"gender" db:"gender"`
	DOB        string `json:"dob" db:"dob"`
	SchoolType string `json:"schooltype" db:"schooltype"`
	SchoolName string `json:"schoolname" db:"schoolname"`
	UID        string `json:"uid" db:"uid"`
}

func NewStudent(fields []byte, u User) error {
	s := make(map[string]interface{})
	err := json.Unmarshal(fields, &s)
	if err != nil {
		return err
	}

	fmt.Println(s["schoolyear"])

	res, err := db.Exec("INSERT INTO students(Gender, DOB, SchoolYear, CurrentLessonID, CurrentChapterName, UID) VALUES($1, $2, $3, $4, $5, $6)",
		s["gender"], s["dob"], s["schoolyear"], CURRENT_LESSON_ID, CURRENT_CHAPTER_NAME, u.UID)
	fmt.Println(res, err)
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

func NewEmployedUser(fields []byte, u User) error {
	e := Employed{}
	if err := json.Unmarshal(fields, &e); err != nil {
		return err
	}

	res, err := db.Exec("INSERT INTO employed(Gender, DOB, CurrentLessonID, CurrentChapterName, UID) VALUES($1, $2, $3, $4, $5)", e.Gender, e.DOB, CURRENT_LESSON_ID, CURRENT_CHAPTER_NAME, u.UID)
	if res != nil {
		return err
	}
	return err
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
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return "", err
	}

	fmt.Println(u)
	fmt.Println(string(fields))

	u.Password = hashedPassword
	result := db.QueryRowx("INSERT INTO users(UID, Firstname, Lastname, Username, Email, Password, Occupation) VALUES($1, $2, $3, $4, $5, $6, $7)",
		u.UID, u.FirstName, u.LastName, u.Username, u.Email, u.Password, u.WhichOccupation).Err()
	if result != nil {
		return "", err
	}

	switch u.WhichOccupation {
	case "student":
		NewStudent(fields, u)
	case "instructor":
		NewInstructor(fields, u)
	case "employed":
		NewEmployedUser(fields, u)
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
	err := db.QueryRowx("SELECT * FROM users WHERE uid=$1", uid).StructScan(&s)
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
