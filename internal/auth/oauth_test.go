package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAuthorizeHandler(t *testing.T) {
	// Prepare query parameters for a valid request.
	params := url.Values{}
	params.Add("client_id", "test_client")
	params.Add("response_type", "code")
	params.Add("state", "test_state")
	// For test, create a dummy code challenge.
	dummyVerifier := "test_code_verifier_which_is_long_enough_for_demo_purposes_12345678"
	codeChallenge := GenerateCodeChallenge(dummyVerifier)
	params.Add("code_challenge", codeChallenge)
	params.Add("code_challenge_method", "S256")

	req := httptest.NewRequest(http.MethodGet, "/authorize?"+params.Encode(), nil)
	rr := httptest.NewRecorder()

	AuthorizeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Check that the response contains a code and correct state.
	if resp["code"] == "" {
		t.Error("Expected an authorization code in response")
	}
	if resp["state"] != "test_state" {
		t.Errorf("Expected state %q, got %q", "test_state", resp["state"])
	}
}

func TestTokenHandler(t *testing.T) {
	// Prepare a valid authorization code and store the expected code challenge.
	codeVerifier := "test_verifier_which_is_long_enough_1234567890abcdef"
	expectedChallenge := GenerateCodeChallenge(codeVerifier)
	authCode := "test_auth_code_token"
	authCodes[authCode] = expectedChallenge

	// Prepare form data.
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode)
	form.Add("code_verifier", codeVerifier)

	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
