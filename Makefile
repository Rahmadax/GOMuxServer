SHELL=/bin/bash
GO        ?= go

HAS_GO       := $(shell command -v go)
HAS_MOCKGEN       := $(shell command -v mockgen)

run-db-migration:
	docker exec -it dbTest bash
	mysql --user=root --password="Password123" dbTest < /migrations/migrations/init.sql

boot:


.PHONY boot:
	make build
	make gen-mocks

build:
	cd Api/cmd;
	go build

gen-mocks:
	cd Api/pkg/guest_list; \
	mockgen -source=guest_list_service.go -destination=mock_guest_list_service.go -package=guest_list; \
	mockgen -source=guest_list_handler.go -destination=mock_guest_list_handler.go -package=guest_list; \
	mockgen -source=guests_service.go -destination=mock_guests_service.go -package=guests; \
	mockgen -source=guests_handler.go -destination=mock_guests_handler.go -package=guests; \
	mockgen -source=seats_service.go -destination=mock_seats_service.go -package=empty_seats; 
