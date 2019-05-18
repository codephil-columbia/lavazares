package routes

import (
	"lavazares/utils"
	"net/http"
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

//CocogameHandler corresponds to GET /game/coco
func CocogameHandler(w http.ResponseWriter, r *http.Request) {
	gameText, err := gameManager.ReturnCocoText()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.SendJSON(gameText, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
