package guest_list

import (
	"encoding/json"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type GuestListService interface {
	getGuestList() (models.GuestList, error)
	addToGuestList(guest models.Guest) error
	removeFromGuestList(guestName string) error
}

type SystemValidator interface {
	IsValidGuestName(name string) bool
	IsValidGuestNumber(accompanyingGuests int) bool
	IsValidTableNumber(tableNumber int) bool
}

type guestListHandler struct {
	config  conf.Configuration
	service GuestListService
	validator SystemValidator
}

func NewGuestListHandler(service GuestListService, validator SystemValidator)  *guestListHandler {
	return &guestListHandler{
		service: service,
		validator: validator,
	}
}

func AddGuestListRoutes(routes conf.RoutesConfig, service GuestListService, validator SystemValidator, router *mux.Router){
	handler := NewGuestListHandler(service, validator)

	router.HandleFunc(routes.GetGuestListUri, handler.getGuestListHandler()).Methods("GET")
	router.HandleFunc(routes.PostGuestListUri, handler.postGuestListHandler()).Methods("POST")
	router.HandleFunc(routes.DeleteGuestListUri, handler.guestListDeleteHandler()).Methods("DELETE")
}

func (glHandler *guestListHandler) getGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestList, err := glHandler.service.getGuestList()
		if err != nil {
			handleErrorResponse(http.StatusBadRequest, err.Error(), w)
		}

		response, _ := json.Marshal(guestList)
		_, _ = w.Write(response)
	}
}

func (glHandler *guestListHandler) postGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, _ := ioutil.ReadAll(r.Body)
		newGuest := models.Guest{}
		err := json.Unmarshal(body, &newGuest)
		if err != nil {
			handleErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}

		newGuest.Name = mux.Vars(r)["name"]
		//if ok := newGuest.validate(glHandler.config.Tables.TableCount); !ok {
		//	handleErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid guest name: %s", newGuest.Name), w)
		//	return
		//}

		err = glHandler.service.addToGuestList(newGuest)
		if err != nil {
			handleErrorResponse(http.StatusBadRequest, err.Error(), w)
		}

		w.WriteHeader(http.StatusCreated)
		response, _ := json.Marshal( models.NameResponse{Name: newGuest.Name})
		_, _ = w.Write(response)
	}
}

func (glHandler *guestListHandler) guestListDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]
		if !glHandler.validator.IsValidGuestName(guestName) {
			handleErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid guest name: %s", guestName), w)
			return
		}

		err := glHandler.service.removeFromGuestList(guestName)
		if err != nil {
			handleErrorResponse(http.StatusInternalServerError, "Something went wrong", w)
		}

		response, _ := json.Marshal(models.NameResponse{Name: guestName})
		_, _ = w.Write(response)
	}
}

func handleErrorResponse(statusCode int, errMessage string, w http.ResponseWriter){
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(errMessage))
}
