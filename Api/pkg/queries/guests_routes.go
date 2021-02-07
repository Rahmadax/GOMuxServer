package queries

const (
	GetGuestList = `
		SELECT name, table, accompanying_guests FROM guests
	`

	InsertGuest = `
		INSERT INTO guests (name, table, accompanying_guests)
		VALUES ($1, $2, $3)
	`

	CountGuestsAtTable = `
		SELECT COUNT(name) + SUM(accompanying_guests) FROM guests WHERE table = %s
	`

	CountPresentGuests = `
		SELECT COUNT(name) + SUM(accompanying_guests) FROM guests
	`

	DeleteGuest = `
		DELETE FROM guests WHERE name = %s
	`
)