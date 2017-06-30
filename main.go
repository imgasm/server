package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/render"
	"github.com/imgasm/server/auth"
	"github.com/imgasm/server/user"
	"net/http"
	"time"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", index)
	r.Mount("/auth", auth.Routes())
	r.Mount("/user", user.Routes())

	http.ListenAndServe(":8080", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}
