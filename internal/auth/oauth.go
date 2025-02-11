package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// Global map to store authorization codes and associated code challenges.
// For demonstration purposes only.
var authCodes = make(map[string]string)

// generateRandomString creates a URL-safe, base64 encoded random string of given byte length.
func generateRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// AuthorizeHandler handles the /authorize endpoint with PKCE.
// Expected query parameters:
// - client_id
// - response_type (must be "code")
// - state
// - code_challenge
// - code_challenge_method (must be "S256")
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	responseType := r.URL.Query().Get("response_type")
	state := r.URL.Query().Get("state")
	codeChallenge := r.URL.Query().Get("code_challenge")
	challengeMethod := r.URL.Query().Get("code_challenge_method")

	// Basic validation.
	if clientID == "" || responseType != "code" || state == "" || codeChallenge == "" || challengeMethod != "S256" {
		http.Error(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	// Generate an authorization code.
	authCode, err := generateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate authorization code", http.StatusInternalServerError)
		return
	}

	// Store the code_challenge for later verification during token exchange.
	authCodes[authCode] = codeChallenge

	// Respond with the auth code and the state.
	resp := map[string]string{
		"code":  authCode,
		"state": state,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// TokenHandler handles the /token endpoint with PKCE validation.
// Expected query parameters (typically via POST form):
// - client_id
// - grant_type (must be "authorization_code")
// - code
// - code_verifier
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form parameters in case POST form is used.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	grantType := r.Form.Get("grant_type")
	code := r.Form.Get("code")
	codeVerifier := r.Form.Get("code_verifier")

	if grantType != "authorization_code" || code == "" || codeVerifier == "" {
		http.Error(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	// Retrieve the stored code challenge.
	storedChallenge, ok := authCodes[code]
	if !ok {
		http.Error(w, "Invalid or expired code", http.StatusBadRequest)
		return
	}

	// Validate using PKCE.
	calculatedChallenge := GenerateCodeChallenge(codeVerifier)
	if calculatedChallenge != storedChallenge {
		http.Error(w, "PKCE validation failed", http.StatusBadRequest)
		return
	}

	// Remove the used authorization code.
	delete(authCodes, code)

	// Respond with an access token.
	resp := map[string]string{
		"access_token": "dummy_access_token",
		"token_type":   "bearer",
		"expires_in":   "3600",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RefreshHandler handles the /refresh endpoint.
// This should verify the refresh token and issue a new access token.
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"access_token": "new_dummy_access_token",
		"token_type":   "bearer",
		"expires_in":   "3600",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
