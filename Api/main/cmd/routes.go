package main

func (s *server)routes(config *Configuration) {

	s.router.HandleFunc("", s.handleGuestListGet()).Methods("GET")
	s.router.HandleFunc("", s.handleGuestListPost()).Methods("POST")
	s.router.HandleFunc("", s.handleGuestListDelete()).Methods("DELETE")

	//s.router.HandleFunc("/invitation/{name}", s.handleInvitationGet()).Methods("GET")
	//
	//s.router.HandleFunc("/guests/{name}", s.handleGuestsGet()).Methods("GET")
	//s.router.HandleFunc("/guests/{name}", s.handleGuestsCreate()).Methods("POST")
}
