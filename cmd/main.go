package main

import (
	"net/http"

	"github.com/zvdy/go-pkce/internal/auth"
)

func main() {
	http.HandleFunc("/authorize", auth.AuthorizeHandler)
	http.HandleFunc("/token", auth.TokenHandler)
	http.HandleFunc("/refresh", auth.RefreshHandler)

	// Optionally, keep your API endpoints from internal/api
	// http.HandleFunc("/api/resource", api.GetAPIResource)

	http.ListenAndServe(":8080", nil)
}
