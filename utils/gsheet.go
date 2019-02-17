package utils

import (
	"os"

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

// GetSheetService returns an editable sheet
func GetSheetService() (*sheets.Service, error) {
	service, err := sheets.New(GetOauthClient())

	if err != nil {
		return nil, err
	}

	return service, nil
}

// GetSheet retrieves the contents of the Google sheet.
func GetSheetData() ([][]interface{}, error) {
	// The sheet service allows us to make requests on a given sheet.
	service, err := GetSheetService()

	if err != nil {
		return nil, err
	}

	// Retrieve sheet data.
	sheetRange, err := service.Spreadsheets.Values.Get(
		os.Getenv("GOOGLE_SPREADSHEET_ID"),
		os.Getenv("GOOGLE_SHEET_RANGE"),
	).ValueRenderOption("FORMATTED_VALUE").Do()

	if err != nil {
		return nil, err
	}

	return sheetRange.Values, nil
}

// ValueToSheet converts a single value to the JSON shape expected by the
// Google Sheets API.
func ValueToSheet(value string) *sheets.ValueRange {
	// Create empty interface.
	row := make([]interface{}, 1)
	values := make([][]interface{}, 1)

	// Fill in interface with given value.
	row[0] = value
	values[0] = row

	return &sheets.ValueRange{
		Values: values,
	}
}
