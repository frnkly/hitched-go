package handlers

// Sheet represents a Google sheet
type Sheet struct {
	name string
}

// GetSheet retrieves the Google sheet with all the guest list information.
func GetSheet() Sheet {
	sheet := Sheet{name: "Hitched Sheet"}

	return sheet
}
