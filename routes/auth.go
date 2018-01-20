package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lavazares/models"
)

//HandleLogin logs in a user
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	req := models.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)

	user, err := models.AutheticateUser(&req)
	if err != nil {
		log.Printf("error logging in: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	sessionID := models.RandomKey()
	session.Values["sessionID"] = sessionID
	session.Values["username"] = user.Username
	session.Save(r, w)

	fmt.Println(sessionID)

	err = models.RedisCache.Set(sessionID, true, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	return
}

//Test is used for test
func Test(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println("in Test")

	val := session.Values["sessionID"]
	sessionID, ok := val.(string)
	if !ok {
		log.Println("error getting settion id")
		http.Error(w, "error", http.StatusForbidden)
	}

	fmt.Println(session)
	fmt.Println(sessionID)
	w.WriteHeader(http.StatusOK)
}

//HandleSignup adds a user to the database
func HandleSignup(w http.ResponseWriter, r *http.Request) {
	newUser := models.User{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading user json: %s", err)
		http.Error(w, "Error Loging in", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &newUser)
	if err != nil {
		log.Printf("error making user: %s", err)
		http.Error(w, "Error loging in", http.StatusBadRequest)
		return
	}

	hashedPassword, err := models.HashPassword(newUser.Password)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		http.Error(w, "Password incorrect", http.StatusBadRequest)
		return
	}

	newUser.Password = hashedPassword

	err = models.InsertUser(&newUser)
	if err != nil {
		log.Printf("error inserting user: %s", err)
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

//HandleLogOut logs out a user and deletes their session
func HandleLogOut(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Printf("Error getting session: %s", err)
		http.Error(w, "Error getting session", http.StatusInternalServerError)
	}
	val := session.Values["sessionID"]
	sessionID, ok := val.(string)
	if !ok {
		log.Printf("Error parsing session id:%s", err)
	}
	models.RedisCache.Del(sessionID)
	w.WriteHeader(http.StatusOK)
	return
}
