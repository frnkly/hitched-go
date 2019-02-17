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

	"google.golang.org/api/sheets/v4"
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

// FirstNameColIndex - Index of "name" column in Google sheet.
const FirstNameColIndex int = 0

// LastNameColIndex - Index of "name" column in Google sheet.
const LastNameColIndex int = 1

// CeremonyInvitationColIndex - Index of "ceremony invitation" column in Google sheet.
const CeremonyInvitationColIndex int = 2

// CeremonyConfirmationColIndex - Index of "ceremony confirmation" column in Google sheet.
const CeremonyConfirmationColIndex int = 3

// ReceptionInvitationColIndex - Index of "reception invitation" column in Google sheet.
const ReceptionInvitationColIndex int = 4

// ReceptionConfirmationColIndex - Index of "reception confirmation" column in Google sheet.
const ReceptionConfirmationColIndex int = 5

// CodeColIndex - Index of "code" column in Google sheet.
const CodeColIndex int = 6

// ContactColIndex - Index of "contact" column in Google sheet.
const ContactColIndex int = 7

// LanguageColIndex - Index of "language" column in Google sheet.
const LanguageColIndex int = 8

// CheckMark is the checkmark charater used in the Google sheet.
const CheckMark string = "✓"

// ---
// Utility functions
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

// GetSheetService returns an editable sheet
func GetSheetService() (*sheets.Service, error) {
	service, err := sheets.New(GetOauthClient())

	if err != nil {
		return nil, err
	}

	return service, nil
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

// ValueToSheet converts a single value to the JSON shape expected by the
// Google Sheets API.
func ValueToSheet(value string) *sheets.ValueRange {
	row := make([]interface{}, 1)
	row[0] = value
	values := make([][]interface{}, 1)
	values[0] = row

	return &sheets.ValueRange{
		Values: values,
	}
}
