package utils

import (
	"errors"
	"log"
	"os"
	"strconv"
)

// ---
// Data structures
// ---

// Guest represents a guest pulled from the Google sheet.
type Guest struct {
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Contact              string `json:"contact"`
	IsAttendingCeremony  int8   `json:"isAttendingCeremony"`
	IsAttendingReception int8   `json:"isAttendingReception"`
	RowIndex             int    `json:"rowIdx"`
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
// Utility functions
// ---

// GetInvitationList returns all invitations, organized by code.
func GetInvitationList() (*InvitationList, error) {
	guestList := InvitationList{Invitations: map[string]*Invitation{}}

	// Get sheet from API.
	sheetData, err := GetSheetData()

	if err != nil {
		return nil, err
	}

	// Populate sheet struct.
	for rowIdx, data := range sheetData {
		// Skip empty rows and guests who don't have a ceremony invite.
		if data[FirstNameColIndex] == "" ||
			len(data) < CeremonyInvitationColIndex ||
			data[CeremonyInvitationColIndex] == "" {
			continue
		}

		// Check invitation code.
		code, isString := data[CodeColIndex].(string)
		if isString == false {
			continue
		} else if len(code) != 4 {
			continue
		}

		// Pull guest data.
		firstNameColData, firstNameTest := data[FirstNameColIndex].(string)
		lastNameColData, lastNameTest := data[LastNameColIndex].(string)
		contactColData := ""
		contactTest := false

		if len(data) > ContactColIndex {
			contactColData, contactTest = data[ContactColIndex].(string)

			if contactTest == false {
				contactColData = ""
			}
		}

		if firstNameTest == false {
			firstNameColData = ""
		}

		if lastNameTest == false {
			lastNameColData = ""
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
			FirstName:            firstNameColData,
			LastName:             lastNameColData,
			Contact:              contactColData,
			IsAttendingCeremony:  isAttending(data[CeremonyConfirmationColIndex]),
			IsAttendingReception: isAttending(data[ReceptionConfirmationColIndex]),
			RowIndex:             rowIdx + 3,
		}

		guestList.Invitations[code].Guests = append(guestList.Invitations[code].Guests, &guest)
	}

	return &guestList, nil
}

// SetAttendence ...
func (g *Guest) SetAttendence(column int, isAttending bool) error {
	// Create sheet range.
	rangeValue := strconv.Itoa(g.RowIndex)

	if column == CeremonyConfirmationColIndex {
		rangeValue = CeremonyConfirmationColLetter + rangeValue
	} else if column == ReceptionConfirmationColIndex {
		rangeValue = ReceptionConfirmationColLetter + rangeValue
	} else {
		return errors.New("invalid attendance column")
	}

	// Google Sheet service.
	service, err := GetSheetService()

	if err != nil {
		return err
	}

	// Attendance value
	attendance := ValueToSheet(CheckMark)

	if isAttending == false {
		attendance = ValueToSheet(CrossMark)
	}

	// Update invitation in sheet.
	response, err := service.Spreadsheets.Values.Update(
		os.Getenv("GOOGLE_SPREADSHEET_ID"),
		rangeValue,
		attendance,
	).ValueInputOption("RAW").Do()

	if err != nil {
		return err
	}

	// log.Println("invitation", invite)
	log.Println("update response", response)

	return nil
}

// --
// Helper functions
// --

func isAttending(data interface{}) int8 {
	attendingData, attendingTest := data.(string)

	if attendingTest == false {
		return 0
	}

	if attendingData == CheckMark {
		return 1
	}

	if attendingData == CrossMark {
		return 2
	}

	return 0
}
