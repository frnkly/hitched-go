package utils

import (
	"errors"
)

// ---
// Data structures
// ---

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
// Helper functions
// ---

// GetInvitationList returns all invitations, organized by code.
func GetInvitationList() (*InvitationList, error) {
	guestList := InvitationList{Invitations: map[string]*Invitation{}}

	// Get sheet from API.
	sheetResponse, err := GetSheetStream()

	if err != nil {
		return nil, errors.New("Google Sheets error: " + err.Error())
	}

	// Parse JSON response.
	sheetValues, err := ParseSheetStream(sheetResponse)

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
