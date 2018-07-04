package routes

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"lavazares/models"
)

//HandleLogin logs in a user
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	req, err := requestToBytes(r.Body)
	if err != nil {
    log.Printf("Error: %v", err)
  }

	_, err = models.AuthenticateUser(req)
	if err != nil {
		log.Printf("User was not found: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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

	newUser, err := models.NewUser(data)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(newUser.GetUserSimpleFields())
	return
}

//HandleLogOut logs out a user and deletes their session
func HandleLogOut(w http.ResponseWriter, r *http.Request) {
	// clientSession, err := store.Get(r, "clientSession")
	// if err != nil {
	// 	log.Printf("Error getting clientSession: %s", err)
	// 	http.Error(w, "Error getting clientSession", http.StatusInternalServerError)
	// }
	// val := clientSession.Values["clientSessionID"]
	// clientSessionID, ok := val.(string)
	// if !ok {
	// 	log.Printf("Error parsing clientSession id:%s", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// err = models.DeleteFromSession(clientSessionID)
	// if err != nil {
	// 	log.Printf("Error deleting from session: %v", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// return
}

// CheckUsernameAvailable checks whether a username has been taken already
func CheckUsernameAvailable(w http.ResponseWriter, r *http.Request) {
  req, err := requestToBytes(r.Body);
	if err != nil {
    log.Printf("Error: %v", err)
  }

  valid, err := models.IsUsernameValid(req);
  if err != nil {
    log.Printf("%v", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  log.Printf("VALID: %v", valid);
  js, _ := json.Marshal(valid)
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
	return
}

func HandleNewPassword(w http.ResponseWriter, r *http.Request) {
  req, _ := requestToBytes(r.Body)
	body := make(map[string]string)
	err := json.Unmarshal(req, &body)

	log.Printf("USERNAME: <%s>", body["username"])
	valid, err := models.IsUsernameValid(req)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
  if !valid {
    log.Printf("VALID? %v | ERROR: %v", valid, err)
    w.WriteHeader(http.StatusOK)
    return
  }
	
	err = models.EditPassword(req)
	if err != nil {
    log.Printf("Password update failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
  }

	/*_, err = models.AuthenticateUser(req)
	if err != nil {
		log.Printf("User was not found: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}*/ // TODO might be optional after auth is finished
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
