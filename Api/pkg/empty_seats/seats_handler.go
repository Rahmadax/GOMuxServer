package empty_seats

import (
	"encoding/json"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
)

type EmptySeatsService interface {
	countEmptySeats() (int, error)
}

type emptySeatsHandler struct {
	config  conf.Configuration
	service EmptySeatsService
}

func NewEmptySeatsHandler(service EmptySeatsService)  *emptySeatsHandler {
	return &emptySeatsHandler {
		service: service,
	}
}

func AddGuestListRoutes(routes conf.RoutesConfig, service EmptySeatsService, router *mux.Router){
	handler := NewEmptySeatsHandler(service)

	router.HandleFunc(routes.CountEmptySeatsUri, handler.countEmptySeatsHandler()).Methods("GET")
}

func (esHandler *emptySeatsHandler) countEmptySeatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentGuestCount, err := esHandler.service.countEmptySeats()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
		}

		response, _ := json.Marshal(models.SeatsEmptyResponse{SeatsEmpty: esHandler.config.Tables.TotalCapacity - presentGuestCount})
		_, _ = w.Write(response)
	}
}