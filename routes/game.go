package api

import (
	"lavazares/utils"
	"net/http"

	"github.com/gorilla/mux"
)


// BoatgameHandler corresponds to GET /game/boatrace 
func BoatgameHandler(w http.ResponseWriter, r *http.Request) {
	gameText, err := gameManager.ReturnBoatText()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.SendJSON(gameText, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
