package auth

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rs/xid"
)

var defaultUserManager *DefaultUserManager
var mockUserStore MockUserStore
var testUser User

type MockUserStore struct{}

func (store MockUserStore) Query(id string) (*User, error) {
	return &testUser, nil
}

func (store MockUserStore) Insert(user *User) error {
	return nil
}

func TestMain(m *testing.M) {
	defaultUserManager = NewDefaultUserManager(mockUserStore)
	testUser = User{
		UID:             xid.New().String(),
		FirstName:       "Test",
		LastName:        "Good",
		Username:        "GoodTester",
		Email:           "tester@gmail.com",
		Password:        "test123",
		Occupation:      "student",
		WhichOccupation: "student",
	}
	ret := m.Run()
	os.Exit(ret)
}

func TestNewUser(t *testing.T) {
	// User missing username
	badUser := User{
		UID:             xid.New().String(),
		FirstName:       "Test",
		LastName:        "Bad",
		Email:           "tester123@gmail.com",
		Password:        "test123",
		Occupation:      "student",
		WhichOccupation: "student",
	}

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
		name     string
		user     User
		expected bool
	}{
		{"Should hash correctly", testUser, false},
		{"Should return error", badUser, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, res := defaultUserManager.hashPassword(tc.user.Password)
			if res != nil && tc.expected == true {
				t.Errorf("expected %t, got %t", tc.expected, res)
			}
		})
	}
}
