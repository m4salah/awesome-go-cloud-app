package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Homepage(mux chi.Router) {
	mux.Get("/home/page", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>My Home Page</h1>")
	})
}
