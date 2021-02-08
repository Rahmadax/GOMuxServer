package queries

const (
	GetGuestList = `
		SELECT guest_name, table_number, accompanying_guests FROM guests ORDER BY table_number asc, guest_name asc
	`

	InsertGuest = `
		INSERT INTO guests (guest_name, table_number, accompanying_guests)
		VALUES (?, ?, ?)
	`

	CountGuestsAtTable = `
		SELECT
		COALESCE(sum(accompanying_guests), 0) + COALESCE(count(guest_name), 0) 
		FROM guests 
		WHERE table_number=?
	`

	CountPresentGuests = `
		SELECT guest_name, accompanying_guests, time_arrived 
		FROM guests 
		WHERE time_arrived IS NOT NULL 
		ORDER BY time_arrived asc, guest_name asc
	`

	DeleteGuest = `
		DELETE FROM guests WHERE guest_name = ?
	`
)