package queries

const (

	// SELECT
	GetGuestList = `
		SELECT guest_name, table_number, accompanying_guests 
		FROM guests 
		ORDER BY table_number asc, guest_name asc
	`

	GetGuestFullDetails = `
		SELECT table_number, accompanying_guests, time_arrived, time_left
		FROM guests 
		WHERE guest_name = ?
	`

	GetPresentGuests = `
		SELECT guest_name, accompanying_guests, time_arrived
		FROM guests 
		WHERE time_arrived IS NOT NULL AND time_left IS NULL
		ORDER BY time_arrived asc, guest_name asc
	`

	// INSERT
	InsertGuest = `
		INSERT INTO guests (guest_name, table_number, accompanying_guests)
		VALUES (?, ?, ?)
	`

	// UPDATE
	UpdateGuestArrives = `
		UPDATE guests
		SET accompanying_guests = ?, time_arrived = ?
		WHERE guest_name=?
	`

	UpdateGuestLeaves = `
		UPDATE guests 
		SET time_left = ?
		WHERE guest_name = ?
	`

	// DELETE
	DeleteFromGuestList = `
		DELETE FROM guests
		WHERE guest_name = ?
	`

	// Count
	CountExpectedGuestsAtTable = `
		SELECT
		COALESCE(sum(accompanying_guests), 0) + COALESCE(count(guest_name), 0) 
		FROM guests 
		WHERE table_number=? AND time_left IS NULL
	`

	CountPresentGuests = `
		SELECT
		COALESCE(sum(accompanying_guests), 0) + COALESCE(count(guest_name), 0) 
		FROM guests
		WHERE time_arrived IS NOT NULL
	`
)
