package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// ---
// Router methods
// ---

// GetRsvp looks-up an invitation code from Google Sheets.
func GetRsvp(writer http.ResponseWriter, request *http.Request) {
	code := chi.URLParam(request, "code")

	// Validate incoming data.
	codeValidator := regexp.MustCompile("^[a-zA-Z0-9-]+$")
	validated := codeValidator.MatchString(code)

	if validated == false {
		render.Render(writer, request, BadRequestError(errors.New("invalid code")))
		return
	}

	// Retrieve invitation.
	invite, err := getInvitationByCode(code)

	if err != nil {
		render.Render(writer, request, NotFoundError(err))
		return
	}

	if err := render.Render(writer, request, Response(invite)); err != nil {
		render.Render(writer, request, RenderingError(err))
		return
	}
}

// Update handles updating guest attendance on the Google spreadsheet.
func Update(writer http.ResponseWriter, request *http.Request) {
	// _ = sheet.Fetch()

	writer.Write([]byte("RSVP'ed"))
}

// ---
// Helper functions
// ---

func getInvitationByCode(code string) (*Invitation, error) {
	// Test code for ceremony invites.
	if code == "test-ceremony" {
		jayne := Guest{Name: "Jayne Mandat", Contact: "jayne.mandat@gmail.com"}
		frank := Guest{Name: "Frank Amankrah", Contact: "frank@frnk.ca"}
		invite := Invitation{
			Code:               code,
			HasCeremonyInvite:  true,
			HasReceptionInvite: false,
			Guests:             []*Guest{&jayne, &frank},
		}

		return &invite, nil
	}

	// Test code for reception invites.
	if code == "test-reception" {
		jayne := Guest{Name: "Jayne Mandat", Contact: "jayne.mandat@gmail.com"}
		frank := Guest{Name: "Frank Amankrah", Contact: "frank@frnk.ca"}
		invite := Invitation{
			Code:               code,
			HasCeremonyInvite:  true,
			HasReceptionInvite: true,
			Guests:             []*Guest{&jayne, &frank},
		}

		return &invite, nil
	}

	// Test code for a large group of guests.
	if code == "test-large" {
		jayne := Guest{Name: "Jayne Mandat", Contact: "jayne.mandat@gmail.com"}
		jasmine := Guest{Name: "Jasmine Mandat", Contact: "jasmine.mandat@gmail.com"}
		judith := Guest{Name: "Judith Mandat", Contact: "judith.mandat@gmail.com"}
		frank := Guest{Name: "Frank Amankrah", Contact: "frank@frnk.ca"}
		invite := Invitation{
			Code:               code,
			HasCeremonyInvite:  true,
			HasReceptionInvite: true,
			Guests:             []*Guest{&jayne, &jasmine, &judith, &frank},
		}

		return &invite, nil
	}

	// Retrieve guest list.
	guestList, err := GetInvitationList()

	if err != nil {
		return nil, err
	}

	// Retrieve invitation.
	if invite, ok := guestList.Invitations[strings.ToUpper(code)]; ok {
		return invite, nil
	}

	return nil, errors.New("code not found")
}

// ---
// Response handlers
// ---

// InvitationResponse is the response payload for the Invitation data model.
//
// In the InvitationResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type InvitationResponse struct {
	*Invitation
}

// Render - renders an InvitationResponse.
func (rd *InvitationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire

	return nil
}

// Response - generates an InvitationResponse response.
func Response(invite *Invitation) render.Renderer {
	return &InvitationResponse{Invitation: invite}
}
