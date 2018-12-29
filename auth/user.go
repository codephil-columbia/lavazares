package auth

import (
	"encoding/json"
	"errors"
	"lavazares/utils"
	"log"
	"os"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jmoiron/sqlx"
)

//User metadata that is stored in the database
type User struct {
	UID             string     `json:"uid"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"-"`
	DeletedAt       *time.Time `json:"-"`
	FirstName       string     `json:"firstName"`
	LastName        string     `json:"lastName"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
	Occupation      string     `json:"occupation"`
	WhichOccupation string     `json:"whichOccupation"`
}

// func (u *User) Unmarshal(byt []byte) error {

// }

const defaultUserManagerLogger = "DefaultUserManager"

// DefaultUserManager handles operations on generic User objects
type DefaultUserManager struct {
	store  UserStore
	logger *log.Logger
}

func NewDefaultUserManager(store UserStore) *DefaultUserManager {
	return &DefaultUserManager{
		store:  store,
		logger: log.New(os.Stdout, defaultUserManagerLogger, log.Lshortfile),
	}
}

// NewUser creates and saves a User
// Before doing so, it will set the UID and hash the password
func (manager *DefaultUserManager) NewUser(args utils.RequestJSON) (*User, error) {
	var user User
	err := json.Unmarshal(args, &user)
	if err != nil {
		return nil, err
	}

	user.UID = xid.New().String()

	hashed, err := manager.hashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashed

	err = manager.store.Insert(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (manager *DefaultUserManager) hashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("string was empty")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// GetUser returns a User by id
func (manager *DefaultUserManager) GetUser(id string) (*User, error) {
	return manager.store.Query(id)
}

type userStore struct {
	db *sqlx.DB
}

type UserStore interface {
	Query(id string) (*User, error)
	Insert(user *User) error
}

func NewUserStore(db *sqlx.DB) *userStore {
	return &userStore{db: db}
}

func (store *userStore) Query(id string) (*User, error) {
	var u User
	err := store.db.QueryRowx("SELECT * FROM Users WHERE id = $1", id).StructScan(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (store *userStore) Insert(user *User) error {
	_, err := store.db.Exec(`INSERT INTO 
		users(UID, Firstname, Lastname, 
		Username, Email, Password, Occupation) 
		VALUES($1, $2, $3, $4, $5, $6, $7)`,
		user.UID, user.FirstName, user.LastName, user.Username,
		user.Email, user.Password, user.WhichOccupation)
	if err != nil {
		return err
	}
	return nil
}