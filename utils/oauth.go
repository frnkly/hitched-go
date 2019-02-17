package utils

// Taken from the Go quickstart guide.
// https://developers.google.com/sheets/api/quickstart/go

import (
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// GetOauthClient uses a service account to generate an access token.
func GetOauthClient() *http.Client {
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
