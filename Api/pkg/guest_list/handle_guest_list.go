package guest_list

import "net/http"

type GuestListHandler struct {

}

func (s *server) handleGuestListGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

//func (s *server) handleGuestListPost() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//	}
//}
//
//func (router *router) handleGuestListDelete() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//	}
//}

func NewGuestListHandler() GuestListHandler {
	return GuestListHandler{}
}