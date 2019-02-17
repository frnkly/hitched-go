package utils

// Taken from the Go quickstart guide.
// https://developers.google.com/sheets/api/quickstart/go

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// ---
// Constants
// ---

const TokenFileName string = "token.json"

// ---
// Utils
// ---

// GetOauthClient uses the service account to generate an access token.
func GetOauthClient() *http.Client {
	log.Println("client email", os.Getenv("GOOGLE_CLIENT_EMAIL"))
	log.Println("client key", os.Getenv("GOOGLE_CLIENT_PRIVATE_KEY"))

	config := &jwt.Config{
		Email:      os.Getenv("GOOGLE_CLIENT_EMAIL"),
		PrivateKey: []byte(os.Getenv("GOOGLE_CLIENT_PRIVATE_KEY")),
		Scopes: []string{
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/drive.file",
			"https://www.googleapis.com/auth/spreadsheets",
		},
		TokenURL: google.JWTTokenURL,
	}

	// Initiate an http.Client.
	return config.Client(oauth2.NoContext)
}

// GetOauthClientFromWeb retrieves a token, saves the token, then returns the generated
// client.
func GetOauthClientFromWeb(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	token, err := getTokenFromFile()

	if err != nil {
		token = getTokenFromWeb(config)
		saveTokenToFile(token)
	}

	return config.Client(context.Background(), token)
}

// ---
// Helper functions
// ---

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)

	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return token
}

// Saves a token to a JSON file.
func saveTokenToFile(token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", TokenFileName)

	tokenFile, err := os.OpenFile(TokenFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer tokenFile.Close()

	json.NewEncoder(tokenFile).Encode(token)
}

// Retrieves a token from a JSON file.
func getTokenFromFile() (*oauth2.Token, error) {

	tokenFile, err := os.Open(TokenFileName)

	if err != nil {
		return nil, err
	}

	defer tokenFile.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(tokenFile).Decode(token)

	return token, err
}
