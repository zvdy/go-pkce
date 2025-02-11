package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zvdy/go-pkce/internal/api"
)

func TestAPIResource(t *testing.T) {
	// Create a new ServeMux and register the same handler as in main.go
	mux := http.NewServeMux()
	mux.HandleFunc("/api/resource", api.GetAPIResource)

	// Start httptest server using the mux
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Create a HTTP GET request for the /api/resource endpoint
	resp, err := http.Get(ts.URL + "/api/resource")
	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}

	// Read the response body using io.ReadAll
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Unmarshal the JSON response
	var result map[string]string
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	// Verify that the "message" field matches the expected output "API Resource"
	expected := "API Resource"
	if msg, ok := result["message"]; !ok || msg != expected {
		t.Fatalf("Expected message %q; got %q", expected, msg)
	}
}
