package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var testBaseRouter *mux.Router
var testUserRouter *mux.Router

func TestMain(m *testing.M) {
	testBaseRouter = mux.NewRouter()
	testUserRouter = testBaseRouter.PathPrefix("/user").Subrouter()
	testBaseRouter.HandleFunc("/edit/password", editPasswordHandler)
}

func TestEditPasswordHandler(t *testing.T) {
	missingUsername := []byte(`{
		"password": "cesar"
	}`)

	missingPassword := []byte(`{
		"username": "cesar"
	}`)

	cases := []struct {
		name     string
		body     []byte
		expected int
	}{
		{"Should not allow missing username", missingUsername, http.StatusBadRequest},
		{"Should not allow missing password", missingPassword, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(
				"POST",
				"/user/edit/password",
				bytes.NewBuffer(tc.body),
			)
			if err != nil {
				t.Errorf("error forming request %v", err)
			}
			testUserRouter.ServeHTTP(rr, req)
			if rr.Code != tc.expected {
				t.Errorf("Expected: [%v], got [%v]", tc.expected, rr.Code)
			}
		})
	}
}
