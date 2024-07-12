package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
)

func main() {
	// Load OAuth 2.0 credentials
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v\n", err)
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/youtube.force-ssl")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v\n", err)
	}
	oauthConfig = config

	getTokenFromWeb(config)

	// Set up HTTP server to handle OAuth callback
	http.HandleFunc("/oauth/callback", handleOAuthCallback)
	fmt.Println("OAuth server started. Go to http://localhost:8080 to initiate OAuth flow...")

	// Listen and serve on localhost:8080
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatalf("Failed to start OAuth server: %v", err)
	}
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for token
	tok, err := exchangeToken(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	// Save token to file
	if err := saveToken("token.json", tok); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save token: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Token successfully saved to token.json. You can close this window.")
}

func exchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	tok, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange token: %w", err)
	}
	return tok, nil
}

func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
}
