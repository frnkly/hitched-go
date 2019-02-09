package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Find
func Find(writer http.ResponseWriter, request *http.Request) {
	code := chi.URLParam(request, "code")
	sheet := GetSheet()

	writer.Write([]byte("RSVP for: " + code + " (" + sheet.name + ")"))
}

// Update handles updating guest attendance on the Google spreadsheet.
func Update(writer http.ResponseWriter, request *http.Request) {
	// _ = sheet.Fetch()

	writer.Write([]byte("RSVP'ed"))
}
