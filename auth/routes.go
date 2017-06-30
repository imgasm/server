package auth

import (
	"github.com/go-chi/chi"
)

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", LoginPOST)
	r.Post("/logout", LogoutPOST)
	r.Post("/register", RegisterPOST)
	return r
}
