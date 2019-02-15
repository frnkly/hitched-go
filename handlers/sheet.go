/**
 * Handles all sheet operations, including authentication, reading and writing
 * to a Google Sheet.
 */
package handlers

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

// Guest represents a guest pulled from the Google sheet.
type Guest struct {
	Name                 string `json:"name"`
	Contact              string `json:"contact"`
	IsAttendingCeremony  bool   `json:"isAttendingCeremony"`
	IsAttendingReception bool   `json:"isAttendingReception"`
}

// Invitation represents a guest's invitation(s), along with their RSVP.
type Invitation struct {
	Code               string   `json:"code"`
	HasCeremonyInvite  bool     `json:"hasCeremonyInvite"`
	HasReceptionInvite bool     `json:"hasReceptionInvite"`
	Guests             []*Guest `json:"guests"`
}

// InvitationList holds the list of all invitations.
type InvitationList struct {
	Invitations map[string]*Invitation
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

// Makes a GET request to the Google Sheets API and returns the response.
// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get
func getSheet() (*http.Response, error) {

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

// Parses a Google Sheets API response into a JSON object.
func parseSheet(stream *http.Response) ([][]string, error) {

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

// GetInvitationList returns all invitations, organized by code.
func GetInvitationList() (*InvitationList, error) {
	guestList := InvitationList{Invitations: map[string]*Invitation{}}

	// Get sheet from API.
	sheetResponse, err := getSheet()

	if err != nil {
		return nil, errors.New("Google Sheets error: " + err.Error())
	}

	// Parse JSON response.
	sheetValues, err := parseSheet(sheetResponse)

	if err != nil {
		return nil, errors.New("JSON parsing error: " + err.Error())
	}

	// Populate sheet struct.
	for _, data := range sheetValues {
		// Skip empty rows and guests who don't have a ceremony invite.
		if data[NameColIndex] == "" || data[CeremonyInvitationColIndex] == "" {
			continue
		}

		// Check invitation code.
		code := data[CodeColIndex]
		if len(code) != 4 {
			continue
		}

		// Create invitation key in guest list.
		if _, ok := guestList.Invitations[code]; !ok {
			invite := Invitation{
				Code:               code,
				HasCeremonyInvite:  true,
				HasReceptionInvite: data[ReceptionInvitationColIndex] == CheckMark,
				Guests:             []*Guest{},
			}

			guestList.Invitations[code] = &invite
		}

		// Add guest to guest list
		guest := Guest{
			Name:                 data[NameColIndex],
			Contact:              data[ContactColIndex],
			IsAttendingCeremony:  data[CeremonyConfirmationColIndex] == CheckMark,
			IsAttendingReception: data[ReceptionConfirmationColIndex] == CheckMark,
		}

		guestList.Invitations[code].Guests = append(guestList.Invitations[code].Guests, &guest)
	}

	return &guestList, nil
}
