package main

import (
    "log"
    "net/http"

    "github.com/zvdy/go-pkce/internal/api"
)

func main() {
    http.HandleFunc("/api/resource", api.GetAPIResource)

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}