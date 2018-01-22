package db

//SaveSession aves user session to session manager
func SaveSession(sessionID, uid string) error {
	return sessionManager.Set(sessionID, uid, 0).Error()
}

//RemoveSession deletes the current session by sessionID
// func RemoveSession(sessionID) error {
// 	return sessionManager.Del(sessionID).Error()
// }
