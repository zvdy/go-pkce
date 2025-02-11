package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorizeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/authorize", nil)
	rr := httptest.NewRecorder()

	AuthorizeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	expectedMessage := "Authorize endpoint - implement client validation and user consent flow here"
	if resp["message"] != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, resp["message"])
	}
}

func TestTokenHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/token", nil)
	rr := httptest.NewRecorder()

	TokenHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if resp["access_token"] != "dummy_access_token" {
		t.Errorf("Expected access_token 'dummy_access_token', got %q", resp["access_token"])
	}
	if resp["token_type"] != "bearer" {
		t.Errorf("Expected token_type 'bearer', got %q", resp["token_type"])
	}
	if resp["expires_in"] != "3600" {
		t.Errorf("Expected expires_in '3600', got %q", resp["expires_in"])
	}
}

func TestRefreshHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	rr := httptest.NewRecorder()

	RefreshHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if resp["access_token"] != "new_dummy_access_token" {
		t.Errorf("Expected access_token 'new_dummy_access_token', got %q", resp["access_token"])
	}
	if resp["token_type"] != "bearer" {
		t.Errorf("Expected token_type 'bearer', got %q", resp["token_type"])
	}
	if resp["expires_in"] != "3600" {
		t.Errorf("Expected expires_in '3600', got %q", resp["expires_in"])
	}
}
