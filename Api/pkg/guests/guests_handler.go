package guests

import (
	"encoding/json"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type GuestsService interface {
	getPresentGuests() (models.PresentGuestList, error)
	guestArrives (int, string) error
	guestLeaves(string) error
}

type SystemValidator interface {
	ValidateGuestName(name string) error
	ValidateArrivingGuest(guestName string, accompanyingGuests int) error
}

type guestsHandler struct {
	config  conf.Configuration
	service GuestsService
	validator SystemValidator
}

func newGuestListHandler(service GuestsService, validator SystemValidator) *guestsHandler {
	return &guestsHandler{
		service: service,
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

func (guestsHandler *guestsHandler) guestArrivesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		updateGuestReq := models.UpdateGuestRequest{}
		err := json.Unmarshal(body, &updateGuestReq)
		if err != nil {
			handleErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		guestName := mux.Vars(r)["name"]
		err = guestsHandler.validator.ValidateArrivingGuest(guestName, updateGuestReq.AccompanyingGuests)
		if err != nil {
			handleErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		err = guestsHandler.service.guestArrives(updateGuestReq.AccompanyingGuests, guestName)
		if err != nil {
			handleErrorResponse(http.StatusBadRequest, err.Error(), w)
		}

		generateHappyResponse(http.StatusOK, models.NameResponse{Name: guestName}, w)
	}
}

func (guestsHandler *guestsHandler) guestLeavesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]
		err := guestsHandler.validator.ValidateGuestName(guestName)
		if err != nil {
			handleErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		err = guestsHandler.service.guestLeaves(guestName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		generateHappyResponse(http.StatusOK, models.NameResponse{Name: guestName}, w)
	}
}

// Invitation
func (guestsHandler *guestsHandler) handleInvitationGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}


func handleErrorResponse(statusCode int, errMessage string, w http.ResponseWriter){
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(errMessage))
}

func generateHappyResponse(statusCode int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	response, _ := json.Marshal(body)
	_, _ = w.Write(response)
}