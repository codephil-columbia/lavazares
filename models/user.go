package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/rs/xid"
)

var (
	FIRST_LESSON_ID    = "d3f9c2a3-1edf-42a6-a24d-3a4ad4683036"
	FIRST_CHAPTER_NAME = "Chapter 0: The Basics"
	FIRST_CHAPTER_ID   = "e6a18785-98c5-41bc-ad98-ec5d3a243d15"
	EMPLOYED           = "employed"
	STUDENT            = "student"
	INSTRUCTOR         = "instructor"
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

func (u *User) GetUserSimpleFields() (map[string]string) {
	return map[string]string {
		"username": u.Username,
		"firstName": u.FirstName,
		"lastName": u.LastName,
		"email": u.Email,
		"uid": u.UID,
	}
}

type Student struct {
	Gender             string `json:"gender" db:"gender"`
	DOB                string `json:"dob" db:"dob"`
	CurrentLessonID    string `json:"currentlessonid" db:"currentlessonid"`
	CurrentChapterID string `json:"currentchapterid" db:"currentchapterid"`
	CurrentChapterName string `db:"currentchaptername"`
	UID                string `db:"uid"`
}

// Have to join w/ Student
type Employed struct {
	Student
	Occupation string `json:"occupation" db:"occupation"`
}

// Have to join w/ Student
type Pupil struct {
	Student
	SchoolYear         string `json:"schoolyear" db:"schoolyear"`
}

type Instructor struct {
	Gender     string `json:"gender" db:"gender"`
	DOB        string `json:"dob" db:"dob"`
	SchoolType string `json:"schooltype" db:"schooltype"`
	SchoolName string `json:"schoolname" db:"schoolname"`
	UID        string `json:"uid" db:"uid"`
}

func NewUser(fields []byte) (*User, error) {

	u := User{}
	err := json.Unmarshal(fields, &u)
	if err != nil {
		return nil, err
	}

	u.UID = xid.New().String()
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashedPassword
	result := db.QueryRowx("INSERT INTO users(UID, Firstname, Lastname, Username, Email, Password, Occupation) VALUES($1, $2, $3, $4, $5, $6, $7)",
		u.UID, u.FirstName, u.LastName, u.Username, u.Email, u.Password, u.WhichOccupation).Err()
	if result != nil {
		return nil, err
	}

	switch u.WhichOccupation {
	case INSTRUCTOR:
		NewInstructor(fields, u)
	case STUDENT:
		NewPupil(fields, u)
	case EMPLOYED:
		NewEmployedStudent(fields, u)
	}

	return &u, err
}

func NewPupil(fields []byte, u User) error {
	e := Pupil{}
	if err := json.Unmarshal(fields, &e); err != nil {
		return err
	}

	res, err := db.Exec("INSERT INTO pupils(schoolyear, uid) VALUES($1, $2)", e.SchoolYear, u.UID)
	if res == nil {
		return err
	}
	return newStudent(fields, u)
}

func NewEmployedStudent(fields []byte, u User) error {
	e := Employed{}
	if err := json.Unmarshal(fields, &e); err != nil {
		return err
	}

	res, err := db.Exec("INSERT INTO employed(occupation, uid) VALUES($1, $2)", e.Occupation, u.UID)
	if res == nil {
		return err
	}

	return newStudent(fields, u)
}

func newStudent(fields []byte, u User) error {
	s := make(map[string]interface{})
	err := json.Unmarshal(fields, &s)
	if err != nil {
		return err
	}
	res, err := db.Exec("INSERT INTO students(Gender, DOB, CurrentLessonID, CurrentChapterName, CurrentChapterID, UID) VALUES($1, $2, $3, $4, $5, $6)",
		s["gender"], s["dob"], FIRST_LESSON_ID, FIRST_CHAPTER_NAME, FIRST_CHAPTER_ID, u.UID)
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

	return &u2, nil
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
