package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ---
// Data structures
// ---

// RawSheet represents the shape of the Google Sheets API response.
type RawSheet struct {
	Range string     `json:"range"`
	Major string     `json:"majorDimension"`
	Rows  [][]string `json:"values"`
}

// ---
// Constants
// ---

// SheetsBaseEndpoint - Base endpoint for Google Sheets API.
const SheetsBaseEndpoint string = "https://sheets.googleapis.com/v4/spreadsheets"

// NameColIndex - Index of "name" column in Google sheet.
const NameColIndex int = 0

// CeremonyInvitationColIndex - Index of "ceremony invitation" column in Google sheet.
const CeremonyInvitationColIndex int = 1

// CeremonyConfirmationColIndex - Index of "ceremony confirmation" column in Google sheet.
const CeremonyConfirmationColIndex int = 2

// ReceptionInvitationColIndex - Index of "reception invitation" column in Google sheet.
const ReceptionInvitationColIndex int = 3

// ReceptionConfirmationColIndex - Index of "reception confirmation" column in Google sheet.
const ReceptionConfirmationColIndex int = 4

// CodeColIndex - Index of "code" column in Google sheet.
const CodeColIndex int = 5

// ContactColIndex - Index of "contact" column in Google sheet.
const ContactColIndex int = 6

// LanguageColIndex - Index of "language" column in Google sheet.
const LanguageColIndex int = 6

// CheckMark is the checkmark charater used in the Google sheet.
const CheckMark string = "✓"

// ---
// Helper functions
// ---

// GetSheetStream makes a GET request to the Google Sheets API and returns
// the response.
// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get
func GetSheetStream() (*http.Response, error) {

	endpointParts := []string{
		SheetsBaseEndpoint,
		"/",
		os.Getenv("GOOGLE_SPREADSHEET_ID"),
		"/values/",
		url.QueryEscape(os.Getenv("GOOGLE_SHEET_RANGE")),
		"?key=",
		os.Getenv("GOOGLE_API_KEY"),
	}

	endpoint := strings.Join(endpointParts, "")

	log.Println("Retrieving Google sheet from:", endpoint)

	stream, err := http.Get(endpoint)

	if err != nil {
		return nil, err
	} else if stream.StatusCode != 200 {
		return nil, errors.New("HTTP Error: " + stream.Status)
	}

	return stream, nil
}

// ParseSheetStream parses a Google Sheets API response into a struct.
func ParseSheetStream(stream *http.Response) ([][]string, error) {

	defer stream.Body.Close()

	// Retrieve JSON string from response body.
	jsonStr, err := ioutil.ReadAll(stream.Body)

	if err != nil {
		return nil, err
	}

	parsedRows := RawSheet{}

	if err := json.Unmarshal(jsonStr, &parsedRows); err != nil {
		return nil, err
	}

	// Return raw values in the sheet.
	return parsedRows.Rows, nil
}
