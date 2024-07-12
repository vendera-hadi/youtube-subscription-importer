package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var oauthConfig *oauth2.Config
var err error

func main() {
	ctx := context.Background()
	// Load OAuth 2.0 credentials
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v\n", err)
		return
	}

	oauthConfig, err = google.ConfigFromJSON(b, youtube.YoutubeForceSslScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v\n", err)
		return
	}

	client, err := getClient(oauthConfig)
	if err != nil {
		// Create a context with cancellation to stop the server after token is saved
		ctx, cancel := context.WithCancel(ctx)
		getTokenFromWeb(oauthConfig)

		// Set up HTTP server to handle OAuth callback
		server := &http.Server{Addr: "localhost:8080"}
		// Set up HTTP server to handle OAuth callback
		http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
			handleOAuthCallback(w, r, cancel)
		})
		fmt.Println("OAuth server started. Copy & paste link above to browser and initiate OAuth flow...")
		// Listen and serve on localhost:8080
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start OAuth server: %v", err)
		}
		// Wait for the context cancellation
		<-ctx.Done()

		// Shutdown the server gracefully
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}

		fmt.Println("Server stopped successfully.")
	}

	importSubscription(ctx, client)
}

func importSubscription(ctx context.Context, client *http.Client) {
	// Create YouTube client
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to create YouTube service: %v\n", err)
		return
	}

	// Open the CSV file
	file, err := os.Open("subscriptions.csv")
	if err != nil {
		fmt.Printf("Unable to open CSV file: %v\n", err)
		return
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Unable to read CSV file: %v\n", err)
		return
	}

	// Subscribe to each channel
	for i, record := range records {
		if i == 0 {
			continue
		}
		channelURL := record[1] // Assuming the Channel URL is in the second column
		channelID, err := extractChannelID(channelURL)
		if err != nil {
			fmt.Printf("Invalid channel URL %s: %v\n", channelURL, err)
			continue
		}
		subscribeToChannel(service, channelID)
	}

	fmt.Println("Subscription process completed.")
}

func extractChannelID(channelURL string) (string, error) {
	parsedURL, err := url.Parse(channelURL)
	if err != nil {
		return "", err
	}

	segments := strings.Split(parsedURL.Path, "/")
	if len(segments) < 3 || segments[1] != "channel" {
		return "", fmt.Errorf("invalid channel URL format")
	}

	return segments[2], nil
}

func subscribeToChannel(service *youtube.Service, channelID string) {
	call := service.Subscriptions.Insert([]string{"snippet"}, &youtube.Subscription{
		Snippet: &youtube.SubscriptionSnippet{
			ResourceId: &youtube.ResourceId{
				Kind:      "youtube#channel",
				ChannelId: channelID,
			},
		},
	})
	_, err := call.Do()
	if err != nil {
		fmt.Printf("Unable to subscribe to channel ID %s: %v\n", channelID, err)
		return
	}
	// fmt.Println(service)
	fmt.Printf("Subscribed to channel ID: %s\n", channelID)
}

func getClient(config *oauth2.Config) (*http.Client, error) {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}
	if tok == nil {
		return nil, fmt.Errorf("token blank")
	}
	return config.Client(context.Background(), tok), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request, cancel context.CancelFunc) {
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

	fmt.Println(w, "Token successfully saved to token.json. Continue to import subscriptions.")

	ctx := context.Background()
	importSubscription(ctx, oauthConfig.Client(ctx, tok))
	cancel() // Cancel the context to stop the server
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
	json.NewEncoder(f).Encode(token)
	return nil
}
