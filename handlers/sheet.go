/**
 * Handles all sheet operations, including authentication, reading and writing
 * to a Google Sheet.
 */
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Sheet represents a Google sheet
type Sheet struct {
	name string
}

const sheetsBaseEndpoint string = "https://sheets.googleapis.com/v4/spreadsheet"

// GetSheet retrieves the Google sheet with all the guest list information.
func GetSheet() Sheet {
	sheet := Sheet{name: "Hitched Sheet"}

	// Retrieve sheet data.
	// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get
	endpoint := []string{
		sheetsBaseEndpoint,
		"/",
		os.Getenv("GOOGLE_SPREADSHEET_ID"),
		"/values/",
		url.QueryEscape(os.Getenv("GOOGLE_SHEET_RANGE")),
		"?key=",
		os.Getenv("GOOGLE_API_KEY"),
	}

	log.Println("Google Sheets endpoint: ", strings.Join(endpoint, ""))

	response, err := http.Get(strings.Join(endpoint, ""))

	if err != nil {
		log.Printf("Google Sheets error: %s\n", err.Error())

		return sheet
	} else if response.StatusCode == 404 {
		log.Println("Sheet not found.")

		return sheet
	} else if response.StatusCode != 200 {
		log.Printf("Response status: %d\n", response.StatusCode)

		return sheet
	}

	// Parse JSON response.
	defer response.Body.Close()

	jsonStr, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Printf("Google Sheets error: %s\n", err.Error())

		return sheet
	}

	log.Println("JSON sheet:", string(jsonStr))

	var decoded map[string]interface{}

	if err := json.Unmarshal(jsonStr, &decoded); err != nil {
		log.Printf("JSON error: %s\n", err.Error())

		return sheet
	}

	fmt.Println(decoded)

	return sheet
}
