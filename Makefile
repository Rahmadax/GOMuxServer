SHELL=/bin/bash
GO        ?= go

HAS_GO       := $(shell command -v go)
HAS_MOCKGEN       := $(shell command -v mockgen)

run-db-migration:
	docker exec -it dbTest bash
	mysql --user=root --password="Password123" dbTest < /migrations/migrations/init.sql

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
	mockgen -source=seats_service.go -destination=mock_seats_service.go -package=empty_seats \
	mockgen -source=seats_handler.go -destination=mock_seats_handler.go -package=empty_seats
