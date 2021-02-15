package guests

import (
	"encoding/json"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
)

type GuestsService interface {
	getPresentGuests() (models.PresentGuestList, error)
	guestArrives(int, string) error
	guestLeaves(string) error
	getInvitation(guestName string) (models.FullGuestDetails, error)
}

type SystemValidator interface {
	ValidateGuestName(name string) error
	ValidateArrivingGuest(guestName string, accompanyingGuests int) error
}

type guestsHandler struct {
	config    conf.Configuration
	service   GuestsService
	validator SystemValidator
}

func newGuestListHandler(service GuestsService, validator SystemValidator) *guestsHandler {
	return &guestsHandler{
		service:   service,
		validator: validator,
	}
}

func AddGuestsRoutes(routes conf.RoutesConfig, service *guestsService, validator SystemValidator, router *mux.Router) {
	handler := newGuestListHandler(service, validator)

	router.HandleFunc(routes.GetGuestsUri, handler.getGuestsHandler()).Methods("GET")
	router.HandleFunc(routes.PutGuestsUri, handler.guestArrivesHandler()).Methods("PUT")
	router.HandleFunc(routes.DeleteGuestsUri, handler.guestLeavesHandler()).Methods("DELETE")

	router.HandleFunc(routes.GetInvitationUri, handler.handleInvitationGet()).Methods("GET")
}

// Guests
func (guestsHandler *guestsHandler) getGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentGuestList, err := guestsHandler.service.getPresentGuests()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		generateHappyResponse(http.StatusOK, presentGuestList, w)
	}
}

func (guestsHandler *guestsHandler) guestLeavesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]
		err := guestsHandler.validator.ValidateGuestName(guestName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = guestsHandler.service.guestLeaves(guestName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		generateHappyResponse(http.StatusOK, models.NameResponse{Name: guestName}, w)
	}
}

func (guestsHandler *guestsHandler) guestArrivesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		updateGuestReq := models.UpdateGuestRequest{}
		err := json.Unmarshal(body, &updateGuestReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := guestsHandler.guestArrives(updateGuestReq, mux.Vars(r)["name"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		generateHappyResponse(http.StatusOK, res, w)

	}
}

func (guestsHandler *guestsHandler) guestArrives(updateGuestReq models.UpdateGuestRequest, guestName string) (models.NameResponse, error) {
	err := guestsHandler.validator.ValidateArrivingGuest(guestName, updateGuestReq.AccompanyingGuests)
	if err != nil {
		return models.NameResponse{}, err
	}

	err = guestsHandler.service.guestArrives(updateGuestReq.AccompanyingGuests, guestName)
	if err != nil {
		return models.NameResponse{}, err
	}

	return models.NameResponse{Name: guestName}, nil
}

// Invitation
func (guestsHandler *guestsHandler) handleInvitationGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type templateStruct struct {
			Name               string
			AccompanyingGuests int
		}

		guestName := mux.Vars(r)["name"]
		err := guestsHandler.validator.ValidateGuestName(guestName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		guestDetails, err := guestsHandler.service.getInvitation(guestName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpl, err := template.ParseFiles("Api/pkg/templates/invite.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_ = tmpl.Execute(w, templateStruct{guestDetails.Name, guestDetails.AccompanyingGuests})
	}
}

func generateHappyResponse(statusCode int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	response, _ := json.Marshal(body)
	_, _ = w.Write(response)
}
