package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/frnkly/hitched/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// --
// Structures
// --

// InvitationResponse is the response payload for the Invitation data model.
//
// In the InvitationResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type InvitationResponse struct {
	*utils.Invitation
}

type RsvpRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

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
		render.Render(writer, request, utils.BadRequestError(errors.New("invalid code")))
		return
	}

	// Retrieve invitation.
	invite, err := getInvitationByCode(code)

	if err != nil {
		render.Render(writer, request, utils.NotFoundError(err))
		return
	}

	if err := render.Render(writer, request, Response(invite)); err != nil {
		render.Render(writer, request, utils.RenderingError(err))
		return
	}
}

func UpdateRsvp(writer http.ResponseWriter, request *http.Request) {
	// Retrieve invitation code.
	code := chi.URLParam(request, "code")
	codeValidator := regexp.MustCompile("^[a-zA-Z0-9-]+$")
	validated := codeValidator.MatchString(code)

	if validated == false {
		render.Render(writer, request, utils.BadRequestError(errors.New("invalid code")))
		return
	}

	// Determine if guest is attending or not.
	action := chi.URLParam(request, "action")
	isAttending := true

	if action == "decline" {
		isAttending = false
	} else if action != "accept" {
		render.Render(writer, request, utils.BadRequestError(errors.New("invalid action")))
		return
	}

	// Retrieve event column index.
	event := chi.URLParam(request, "event")
	eventIndex := 0

	if event == "ceremony" {
		eventIndex = utils.CeremonyConfirmationColIndex
	} else if event == "reception" {
		eventIndex = utils.ReceptionConfirmationColIndex
	} else {
		render.Render(writer, request, utils.BadRequestError(errors.New("invalid event")))
		return
	}

	// Retrieve guest details.
	var target RsvpRequest
	err := json.NewDecoder(request.Body).Decode(&target)

	if target.FirstName == "" && target.LastName == "" {
		render.Render(writer, request, utils.BadRequestError(errors.New("invalid guest details")))
		return
	}

	// Retrieve invitation.
	invite, err := getInvitationByCode(code)

	if err != nil {
		render.Render(writer, request, utils.NotFoundError(err))
		return
	}

	// Find guest.
	for _, candidate := range invite.Guests {
		if candidate.FirstName == target.FirstName && candidate.LastName == target.LastName {
			err = candidate.SetAttendence(eventIndex, isAttending)
			if err != nil {
				render.Render(writer, request, utils.NotFoundError(err))
				return
			} else {
				writer.Write([]byte(""))
				return
			}
		}
	}

	render.Render(writer, request, utils.NotFoundError(errors.New("guest not found")))
}

// ---
// Helper functions
// ---

func getInvitationByCode(code string) (*utils.Invitation, error) {
	// Test code for ceremony invites.
	if code == "test-ceremony" {
		jayne := utils.Guest{
			FirstName:           "Jayne",
			LastName:            "Mandat",
			Contact:             "jayne.mandat@gmail.com",
			IsAttendingCeremony: utils.IsAttendingEvent,
		}

		frank := utils.Guest{
			FirstName: "Frank",
			LastName:  "Amankrah",
			Contact:   "frank@frnk.ca",
		}

		invite := utils.Invitation{
			Code:               code,
			HasCeremonyInvite:  true,
			HasReceptionInvite: false,
			Guests:             []*utils.Guest{&jayne, &frank},
		}

		return &invite, nil
	}

	// Test code for reception invites.
	if code == "test-reception" {
		jayne := utils.Guest{
			FirstName:            "Jayne",
			LastName:             "Mandat",
			Contact:              "jayne.mandat@gmail.com",
			IsAttendingCeremony:  utils.IsAttendingEvent,
			IsAttendingReception: utils.IsAttendingEvent,
		}

		frank := utils.Guest{
			FirstName: "Frank",
			LastName:  "Amankrah",
			Contact:   "frank@frnk.ca",
		}

		invite := utils.Invitation{
			Code:               code,
			HasCeremonyInvite:  true,
			HasReceptionInvite: true,
			Guests:             []*utils.Guest{&jayne, &frank},
		}

		return &invite, nil
	}

	// Test code for a large group of guests.
	if code == "test-large" {
		jayne := utils.Guest{
			FirstName:           "Jayne",
			LastName:            "Mandat",
			Contact:             "jayne.mandat@gmail.com",
			IsAttendingCeremony: utils.IsAttendingEvent,
		}

		jasmine := utils.Guest{
			FirstName: "Jasmine",
			LastName:  "Mandat",
			Contact:   "jasmine.mandat@gmail.com",
		}

		judith := utils.Guest{
			FirstName: "Judith",
			LastName:  "Mandat",
			Contact:   "judith.mandat@gmail.com",
		}

		frank := utils.Guest{
			FirstName: "Frank",
			LastName:  "Amankrah",
			Contact:   "frank@frnk.ca",
		}

		invite := utils.Invitation{
			Code:               code,
			HasCeremonyInvite:  true,
			HasReceptionInvite: true,
			Guests:             []*utils.Guest{&jayne, &jasmine, &judith, &frank},
		}

		return &invite, nil
	}

	// Retrieve guest list.
	guestList, err := utils.GetInvitationList()

	if err != nil {
		return nil, err
	}

	// Retrieve invitation.
	if invite, ok := guestList.Invitations[strings.ToUpper(code)]; ok {
		return invite, nil
	}

	return nil, errors.New("code not found")
}

func updateInvitation( /*invite *utils.Invitation*/ ) error {
	service, err := utils.GetSheetService()

	if err != nil {
		return err
	}

	// Update invitation in sheet.
	updateRange := "H1"
	response, err := service.Spreadsheets.Values.Update(
		os.Getenv("GOOGLE_SPREADSHEET_ID"),
		updateRange,
		utils.ValueToSheet("✓"),
	).ValueInputOption("RAW").Do()

	if err != nil {
		return err
	}

	// log.Println("invitation", invite)
	log.Println("update response", response)

	return errors.New("Updating invitations not implemented.")
}

// ---
// Response handlers
// ---

// Render - renders an InvitationResponse.
func (rd *InvitationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire

	return nil
}

// Response - generates an InvitationResponse response.
func Response(invite *utils.Invitation) render.Renderer {
	return &InvitationResponse{Invitation: invite}
}
