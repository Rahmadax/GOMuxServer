package models

// Models
type PresentGuestList struct {
	Guests []PresentGuest `json:"guests"`
}

type PresentGuest struct {
	Name               string  `json:"name" db:"guest_name"`
	AccompanyingGuests int     `json:"accompanying_guests" db:"accompanying_guests"`
	TimeArrived        string  `json:"time_arrived" db:"time_arrived"`
	TimeLeft           *string `json:",omitempty" db:"time_left"`
}

type UpdateGuestRequest struct {
	AccompanyingGuests int `json:"accompanying_guests" db:"accompanying_guests"`
}

type GuestList struct {
	Guests []Guest `json:"guests"`
}

type Guest struct {
	Name               string `json:"name" db:"guest_name"`
	Table              int    `json:"table" db:"table_number"`
	AccompanyingGuests int    `json:"accompanying_guests" db:"accompanying_guests"`
}

type FullGuestDetails struct {
	Name               string  `db:"guest_name"`
	Table              int     `db:"table_number"`
	AccompanyingGuests int     `db:"accompanying_guests"`
	TimeArrived        *string `db:"time_arrived"`
	TimeLeft           *string `db:"time_left"`
}

type NameResponse struct {
	Name string `json:"name"`
}

type SeatsEmptyResponse struct {
	SeatsEmpty int `json:"seats_empty"`
}