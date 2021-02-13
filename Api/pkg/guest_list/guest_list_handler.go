package guest_list

import (
	"encoding/json"
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
	ValidateGuestName(name string) error
	ValidateNewGuest(newGuest models.Guest, name string) error
}

func AddGuestListRoutes(routes conf.RoutesConfig, service GuestListService, validator SystemValidator, router *mux.Router) {
	handler := NewGuestListHandler(service, validator)

	router.HandleFunc(routes.GetGuestListUri, handler.getGuestListHandler()).Methods("GET")
	router.HandleFunc(routes.PostGuestListUri, handler.postGuestListHandler()).Methods("POST")
	router.HandleFunc(routes.DeleteGuestListUri, handler.guestListDeleteHandler()).Methods("DELETE")
}

// Http Handlers + Marshallers
func (glHandler *guestListHandler) getGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := glHandler.getGuestList()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(res)
		_, _ = w.Write(response)
	}
}

func (glHandler *guestListHandler) postGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		newGuest := models.Guest{}
		err := json.Unmarshal(body, &newGuest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := glHandler.postGuestList(newGuest, mux.Vars(r)["name"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		response, _ := json.Marshal(res)
		_, _ = w.Write(response)
	}
}

func (glHandler *guestListHandler) guestListDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]
		err := glHandler.guestListDelete(guestName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(models.NameResponse{Name: guestName})
		_, _ = w.Write(response)
	}
}


// Validation + Service calls
func (glHandler *guestListHandler) getGuestList() (models.GuestList, error) {
	return glHandler.service.getGuestList()
}

func (glHandler *guestListHandler) postGuestList(newGuest models.Guest, name string) (models.NameResponse, error) {
	err := glHandler.validator.ValidateNewGuest(newGuest, name)
	if err != nil {
		return models.NameResponse{}, err
	}
	newGuest.Name = name

	err = glHandler.service.addToGuestList(newGuest)
	if err != nil {
		return models.NameResponse{}, err
	}

	return models.NameResponse{Name: newGuest.Name}, nil
}

func (glHandler *guestListHandler) guestListDelete(guestName string) error {
	err := glHandler.validator.ValidateGuestName(guestName)
	if err != nil {
		return err
	}

	err = glHandler.service.removeFromGuestList(guestName)
	if err != nil {
		return err
	}

	return nil
}

// Init
type guestListHandler struct {
	config    conf.Configuration
	service   GuestListService
	validator SystemValidator
}

func NewGuestListHandler(service GuestListService, validator SystemValidator) *guestListHandler {
	return &guestListHandler{
		service:   service,
		validator: validator,
	}
}
