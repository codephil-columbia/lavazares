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

// TODO: Add required struct tags

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

const defaultUserManagerLogger = "DefaultUserManager"

// DefaultUserManager handles operations on generic User objects
type DefaultUserManager struct {
	store  UserStore
	logger *log.Logger
}

// NewDefaultUserManager performs operations on User objects.
// It is also responsible for User authentication
func NewDefaultUserManager(store UserStore) *DefaultUserManager {
	return &DefaultUserManager{
		store:  store,
		logger: log.New(os.Stdout, defaultUserManagerLogger, log.Lshortfile),
	}
}

// EditPassword changes and reshashes a Users password
func (manager *DefaultUserManager) EditPassword(username, password string) error {
	user, err := manager.store.QueryByUsername(username)
	if err != nil {
		return err
	}
	newHashed, err := manager.hashPassword(user.Password)
	if err != nil {
		return err
	}
	return manager.store.UpdateUserByUsername(user.Username, "password", newHashed)
}

// Authenticate authenticates a user by checking whether there exists a username
// password pair thats in the db. If the password and hash match, error returned is nil.
func (manager *DefaultUserManager) Authenticate(username, password string) error {
	user, err := manager.store.QueryByUsername(username)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (manager *DefaultUserManager) removeUserByUsername(username string) error {
	return manager.store.DeleteByUsername(username)
}

// NewUser creates and saves a User
// Before doing so, it will set the UID and hash the password
// TODO: Maybe change this to take in the user post unmarsheling?
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

// IsUsernameValid checks to see whether a username is valid. Validity is defined
// as no two usernames should be the same.
func (manager *DefaultUserManager) IsUsernameValid(username string) bool {
	_, err := manager.store.QueryByUsername(username)
	return err != nil
}

func (manager *DefaultUserManager) hashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password was empty")
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

// GetUserByUsername returns a user given a specific username
func (manager *DefaultUserManager) GetUserByUsername(username string) (*User, error) {
	return manager.store.QueryByUsername(username)
}

// userStore satisfies UserStore interface
type userStore struct {
	db *sqlx.DB
}

// UserStore defines User db operations
type UserStore interface {
	Query(id string) (*User, error)
	QueryByUsername(username string) (*User, error)
	Insert(user *User) error
	DeleteByUsername(username string) error
	UpdateUserByUsername(username, field, value string) error
}

// NewUserStore creates a new generic userStore with the
// given db param.
func NewUserStore(db *sqlx.DB) UserStore {
	return &userStore{db: db}
}

func (store *userStore) UpdateUserByUsername(username, field, value string) error {
	_, err := store.db.Exec("UPDATE users SET $1 = $2 WHERE username = $3", field, value, username)
	if err != nil {
		return err
	}
	return nil
}

func (store *userStore) DeleteByUsername(username string) error {
	_, err := store.db.Exec("DELETE FROM Users WHERE username = $1", username)
	if err != nil {
		return err
	}
	return nil
}

func (store *userStore) QueryByUsername(username string) (*User, error) {
	var user User
	err := store.db.QueryRowx("SELECT * FROM Users WHERE username = $1", username).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (store *userStore) Query(id string) (*User, error) {
	var u User
	err := store.db.QueryRowx("SELECT * FROM Users WHERE uid = $1", id).StructScan(&u)
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
