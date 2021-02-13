.PHONY: build-app
	go mod tidy
	go mod download

.PHONY: start-local-db
	docker-compose -f

#mockgen -source=guest_list_service.go -destination=mock_guest_list_service.go -package=guest_list
#mockgen -source=guest_list_handler.go -destination=mock_guest_list_handler.go -package=guest_list

#mockgen -source=seats_service.go -destination=mock_seats_service.go -package=empty_seats
#mockgen -source=seats_handler.go -destination=mock_seats_handler.go -package=empty_seats

