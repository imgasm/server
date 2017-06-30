package auth

import (
	// "github.com/go-chi/chi/render"
	"github.com/imgasm/server/session"
	"net/http"
)

// curl -i -X POST http://localhost:8080/auth/logout
func LogoutPOST(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(r)
	if sess.Values["username"] == nil {
		http.Error(w, "You are not logged in.", 409)
		return
	}
	// session.ClearSessionByKey("username", sess)
	session.ClearAllSessionValues(sess)
	sess.AddFlash("Successfully logged out.")
	sess.Save(r, w)
	w.Write([]byte("Successfully logged out."))
}
