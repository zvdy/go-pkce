package api

import (
    "encoding/json"
    "net/http"
)

func GetAPIResource(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{"message": "API Resource"}
    json.NewEncoder(w).Encode(response)
}