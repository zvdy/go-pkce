# OAuth 2.0 flow with PKCE

1. **Initial Request (/authorize):**  
   - The client generates:
     - A random state (min. 28 characters)  
     - A code verifier (min. 64 characters)  
     - A code challenge by hashing the code verifier (typically using SHA256)
   - Both the state and code verifier are saved locally (e.g., in Local Storage).
   - The client sends a request to the `/authorize` endpoint with the following parameters:
     - client_id  
     - response_type (set to "code")  
     - state (the generated state)  
     - code_challenge (the generated challenge)  
     - code_challenge_method (e.g., "S256")
   - The authorization server validates the incoming request and ensures the code_challenge is valid.

2. **Authorization Response:**  
   - If validated, the authorization server responds by redirecting the user's browser to a callback URL.  
   - The callback URL contains:
     - An authorization code  
     - The state that was originally sent.
   - The client (browser) verifies that the returned state matches the previously stored state.

3. **Token Exchange (/token):**  
   - Once the state is confirmed, the client sends a request to the `/token` endpoint.  
   - This request includes:
     - client_id  
     - grant_type (set to "authorization_code")  
     - the authorization code from the callback  
     - the original code_verifier
   - On the server side, the authorization server regenerates the code challenge using the provided code_verifier (using SHA256) and compares it with the stored challenge value from the initial request.
   - If they match, the server issues an access token along with related parameters (token type, expiry, etc.).

4. **Token Refresh (/refresh):**  
   - When the access token expires, the client can request a new token by calling the `/refresh` endpoint with a valid refresh token.
   - The server validates the refresh token and issues a new access token.

> This flow ensures that even if an attacker intercepts the authorization code, they cannot exchange it for a valid access token without also possessing the original code_verifier.