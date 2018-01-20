package routes

import (
	"log"
	"net/http"

	"github.com/lavazares/models"
)

//AuthMiddleware checks to see if session is valid, else throws error
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session")
		if err != nil {
			log.Printf("error getting session: %s\n", err)
			http.Error(w, "Session Invalid", http.StatusForbidden)
		}

		val := session.Values["sessionID"]
		sessionID, ok := val.(string)
		if !ok || sessionID == "" {
			log.Printf("err or session id was empty: %s", err)
			http.Error(w, "Session ID invalid", http.StatusForbidden)
			return
		}

		exists := models.RedisCache.Exists(sessionID)
		if exists.Val() == 0 {
			log.Printf("Session id does not exist:")
			http.Error(w, "Session ID invalid", http.StatusForbidden)
			return
		}

		log.Printf("user is authenticated")
		next.ServeHTTP(w, r)
	})
}
