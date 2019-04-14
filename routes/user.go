package routes

import (
	"encoding/json"
	"errors"
	"lavazares/utils"
	"net/http"
)

func editPasswordHandler(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ReadBodyToMap(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username, ok := data["username"]
	if !ok {
		http.Error(w, errors.New("Missing username").Error(), http.StatusBadRequest)
		return
	}
	password, ok := data["password"]
	if !ok {
		http.Error(w, errors.New("Missing password").Error(), http.StatusBadRequest)
		return
	}

	err = userManager.EditPassword(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func newUserHandler(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ReadBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = userManager.NewUser(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ReadBodyToMap(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username, ok := data["username"]
	if !ok {
		http.Error(w, "Username was empty", http.StatusBadRequest)
		return
	}

	password, ok := data["password"]
	if !ok {
		http.Error(w, "Password was empty", http.StatusBadRequest)
		return
	}

	user, err := userManager.Authenticate(username, password)
	if err != nil {
		// Make sure not send stack errors
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
