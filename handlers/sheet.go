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

// Guest represents a guest pulled from the Google sheet.
type Guest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Address              string `json:"address"`
	IsAttendingCeremony  bool   `json:"isAttendingCeremony"`
	IsAttendingReception bool   `json:"isAttendingReceptiom"`
}

// Invitation represents a guest's invitation(s), along with their RSVP.
type Invitation struct {
	Code               string   `json:"code"`
	HasCeremonyInvite  bool     `json:"hasCeremonyInvite"`
	HasReceptionInvite bool     `json:"hasReceptionInvite"`
	Guests             []*Guest `json:"guests"`
}

// Sheet holds the list of all invitations.
type Sheet struct {
	invitations []*Invitation
}

const sheetsBaseEndpoint string = "https://sheets.googleapis.com/v4/spreadsheets"

// Makes a GET request to the Google Sheets API and returns the response.
// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get
func getSheet() (*http.Response, error) {

	endpointParts := []string{
		sheetsBaseEndpoint,
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

	// Parse JSON string into a struct.
	var decoded map[string]interface{}

	if err := json.Unmarshal(jsonStr, &decoded); err != nil {
		return nil, err
	}

	// Retrieve column values from struct.
	type ValuesType struct {
		values []string
	}

	rows := decoded["values"].([]interface{})
	parsedRows := make([][]string, len(rows))

	for i := range rows {
		row := rows[i].([]interface{})
		parsedRow := make([]string, len(row))

		for j := range row {
			col := row[j].(string)
			parsedRow[j] := col
		}

		parsedRows[i] = parsedRow
	}

	log.Println("parsedRows", parsedRows)

	// Return raw values in the sheet.
	return parsedRows, nil
}

// GetSheet retrieves the Google sheet with all the guest list information.
func GetSheet() Sheet {
	sheet := Sheet{}

	// Get sheet from API.
	sheetResponse, err := getSheet()

	if err != nil {
		log.Printf("Google Sheets error: %s\n", err.Error())

		return sheet
	}

	// Parse JSON response.
	sheetValues, err := parseSheet(sheetResponse)

	if err != nil {
		log.Printf("JSON parsing error: %s\n", err.Error())

		return sheet
	}

	// Populate sheet struct.
	for _, data := range sheetValues {
		log.Println(data)
	}

	return sheet
}
