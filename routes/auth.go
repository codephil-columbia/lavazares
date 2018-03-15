package routes

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lavazares/models"
)

//HandleLogin logs in a user
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	req, err := requestToBytes(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	u, err := models.AutheticateUser(req)
	if err != nil {
		log.Printf("User was not found: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	serverSession := models.NewUserSession(u.UID)

	clientSession, err := store.Get(r, "session")
	if err != nil {
		log.Printf("Could not get client session: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	clientSession.Values["sessionID"] = serverSession.SessionID
	fmt.Println(serverSession.SessionID)
	err = models.SetToSession(serverSession.SessionID, serverSession.UserID)
	if err != nil {
		log.Printf("Could not commit to server side session: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

//HandleSignup adds a user to the database
func HandleSignup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := models.NewUser(data)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "session")
	session.Values["sessionID"] = id
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	return
}

//HandleLogOut logs out a user and deletes their session
func HandleLogOut(w http.ResponseWriter, r *http.Request) {
	clientSession, err := store.Get(r, "clientSession")
	if err != nil {
		log.Printf("Error getting clientSession: %s", err)
		http.Error(w, "Error getting clientSession", http.StatusInternalServerError)
	}
	val := clientSession.Values["clientSessionID"]
	clientSessionID, ok := val.(string)
	if !ok {
		log.Printf("Error parsing clientSession id:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = models.DeleteFromSession(clientSessionID)
	if err != nil {
		log.Printf("Error deleting from session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func requestToBytes(body io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
