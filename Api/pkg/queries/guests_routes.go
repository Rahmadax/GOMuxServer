package queries

const (
	GetGuestList = `
		SELECT guest_name, table_number, accompanying_guests 
		FROM guests 
		ORDER BY table_number asc, guest_name asc
	`

	InsertGuest = `
		INSERT INTO guests (guest_name, table_number, accompanying_guests)
		VALUES (?, ?, ?)
	`

	GetGuest = `
		SELECT table_number, accompanying_guests
		FROM guests 
		WHERE guest_name = ?
	`

	DeleteAll = `
		DELETE FROM guests WHERE guest_name is not null
	`

	CountExpectedGuestsAtTable = `
		SELECT
		COALESCE(sum(accompanying_guests), 0) + COALESCE(count(guest_name), 0) 
		FROM guests 
		WHERE table_number=?
	`

	UpdateArrivedGuest = `
		UPDATE guests
		SET accompanying_guests = ?, time_arrived = ?
		WHERE guest_name=?
	`

	CountPresentGuests = `
		SELECT
		COALESCE(sum(accompanying_guests), 0) + COALESCE(count(guest_name), 0) 
		FROM guests
		WHERE time_arrived IS NOT NULL
	`

	GetPresentGuests = `
		SELECT guest_name, accompanying_guests, time_arrived 
		FROM guests 
		WHERE time_arrived IS NOT NULL 
		ORDER BY time_arrived asc, guest_name asc
	`

	DeleteGuest = `
		DELETE
		FROM guests 
		WHERE guest_name = ?
	`
)