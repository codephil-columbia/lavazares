package routes

import (
	"github.com/gorilla/sessions"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

var store = sessions.NewCookieStore([]byte("change-this-at-some-point"))
var db = db.DB
