package utils

// ---
// Data structures
// ---

// Guest represents a guest pulled from the Google sheet.
type Guest struct {
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
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
	sheetData, err := GetSheetData()

	if err != nil {
		return nil, err
	}

	// Populate sheet struct.
	for _, data := range sheetData {
		// Skip empty rows and guests who don't have a ceremony invite.
		if data[FirstNameColIndex] == "" || data[CeremonyInvitationColIndex] == "" {
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
			IsAttendingCeremony:  data[CeremonyConfirmationColIndex] == CheckMark,
			IsAttendingReception: data[ReceptionConfirmationColIndex] == CheckMark,
		}

		guestList.Invitations[code].Guests = append(guestList.Invitations[code].Guests, &guest)
	}

	return &guestList, nil
}
