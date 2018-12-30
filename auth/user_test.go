package auth

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

var (
	testDB             *sqlx.DB
	defaultUserManager *DefaultUserManager
	mockUserStore      MockUserStore
	store              UserStore
	uid                = xid.New().String()
	testUser           = User{
		UID:             uid,
		FirstName:       "Test",
		LastName:        "Good",
		Username:        "GoodTester",
		Email:           "tester@gmail.com",
		Password:        "test123",
		Occupation:      "student",
		WhichOccupation: "student",
	}
	badUser = User{
		UID:             xid.New().String(),
		FirstName:       "Test",
		LastName:        "Bad",
		Email:           "tester123@gmail.com",
		Password:        "test123",
		Occupation:      "student",
		WhichOccupation: "student",
	}
)

const (
	conn = "port=5432 host=localhost sslmode=disable user=postgres dbname=postgres password=postgres"
)

type MockUserStore struct{}

func (store MockUserStore) Query(id string) (*User, error) {
	return &testUser, nil
}

func (store MockUserStore) Insert(user *User) error {
	return nil
}

func (store MockUserStore) QueryByUsername(username string) (*User, error) {
	return &testUser, nil
}

func (store MockUserStore) DeleteByUsername(username string) error {
	return nil
}

func (store MockUserStore) UpdateUserByUsername(username, field, value string) error {
	return nil
}

func addTestData() {
	_, err := testDB.Exec(`INSERT INTO 
		users(UID, Firstname, Lastname, 
		Username, Email, Password, Occupation) 
		VALUES($1, $2, $3, $4, $5, $6, $7)`,
		testUser.UID, testUser.FirstName, testUser.LastName, testUser.Username,
		testUser.Email, testUser.Password, testUser.WhichOccupation)
	if err != nil {
		log.Panicf("Could not add test user: %v\n", err.Error())
	}
}

func removeTestData() {
	_, err := testDB.Exec("DELETE FROM Users WHERE uid = $1", uid)
	if err != nil {
		log.Panicln("Could not delete tester user")
	}
}

func TestMain(m *testing.M) {
	defaultUserManager = NewDefaultUserManager(mockUserStore)

	var err error
	testDB, err = sqlx.Open("postgres", conn)
	if err != nil {
		log.Panicf("could not connect to test db: %v\n", err)
	}
	defer testDB.Close()

	store = NewUserStore(testDB)

	ret := m.Run()
	os.Exit(ret)
}

func TestNewUser(t *testing.T) {
	cases := []struct {
		name     string
		user     User
		expected bool
	}{
		{"Insert good User", testUser, true},
		{"Attempt to insert incomplete user", badUser, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := json.Marshal(tc.user)
			if err != nil {
				t.Error("Error marshalling test case")
			}
			_, result := defaultUserManager.NewUser(request)
			if result != nil {
				t.Errorf("expected %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	badUser := User{
		Password: "",
	}
	cases := []struct {
		name        string
		user        User
		expectedErr bool
	}{
		{"Should hash correctly", testUser, true},
		{"Should return error", badUser, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, res := defaultUserManager.hashPassword(tc.user.Password)
			if res != nil && tc.expectedErr {
				t.Errorf("expected %t, got %t", tc.expectedErr, res)
			}
		})
	}
}

// Store Tests
func TestQuery(t *testing.T) {
	cases := []struct {
		name        string
		id          string
		expectedErr bool
	}{
		{"Query with existing user", uid, false},
		{"Query with non existing user", "4321", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			addTestData()
			defer removeTestData()

			user, err := store.Query(tc.id)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("got %v, expected %v", err.Error(), tc.expectedErr)
				}
			} else {
				if user.UID != uid {
					t.Errorf("got user with incorrect id, got %v, expected %v", user.UID, tc.id)
				}
			}
		})
	}
}

func TestQueryByUsername(t *testing.T) {
	cases := []struct {
		name        string
		username    string
		expectedErr bool
	}{
		{"Query with existing user with given username", testUser.Username, false},
		{"Query with nonexistent username", "lilpump", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			addTestData()
			defer removeTestData()

			user, err := store.QueryByUsername(tc.username)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("got [%v], was expecting error to be [%v]", err.Error(), tc.expectedErr)
				}
			} else {
				if user.Username != tc.username {
					t.Errorf("expected username to be [%v], got [%v]", tc.username, user.Username)
				}
			}
		})
	}
}

func TestInsert(t *testing.T) {
	newUser := User{
		UID:             xid.New().String(),
		FirstName:       "Test",
		LastName:        "Good",
		Username:        "TesterGood",
		Email:           "tester@gmail.com",
		Password:        "test123",
		Occupation:      "student",
		WhichOccupation: "student",
	}

	cases := []struct {
		name        string
		user        User
		expectedErr bool
	}{
		{"Insert User with all fields complete", newUser, false},
		{"Insert User with some required fields missing", badUser, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if !tc.expectedErr {
					testDB.Exec("DELETE FROM Users WHERE uid = $1", tc.user.UID)
				}
			}()
			err := store.Insert(&tc.user)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("got [%v], expected err to be [%v]", err.Error(), tc.expectedErr)
				}
			}
		})
	}
}

func TestDeleteByUsername(t *testing.T) {
	cases := []struct {
		name        string
		user        User
		expectedErr bool
	}{
		{"Delete existing User", testUser, false},
		{"Delete nonexisting User", User{Username: "Sang"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			addTestData()
			defer func() {
				if tc.expectedErr {
					removeTestData()
				}
			}()

			err := store.DeleteByUsername(tc.user.Username)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("got [%v], expected [%v]", err, tc.expectedErr)
				}
			}
		})
	}
}

func TestUpdateUserByUsername(t *testing.T) {
	cases := []struct {
		name        string
		user        User
		field       string
		value       string
		expectedErr bool
	}{
		{"Successfully updates field in user", testUser, "username", "lilpump", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			addTestData()
			defer removeTestData()

			err := store.UpdateUserByUsername(tc.user.Username, tc.field, tc.value)
			if err != nil {
				if !tc.expectedErr {
					log.Printf("")
				}
			}
		})
	}
}
