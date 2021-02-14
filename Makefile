SHELL=/bin/bash
GO        ?= go

HAS_GO       := $(shell command -v go)
HAS_MOCKGEN       := $(shell command -v mockgen)

.PHONY boot:
	make build
	make gen-mocks

build:
	cd Api/cmd;

gen-mocks:
	cd Api/pkg/guest_list; \
	mockgen -source=guest_list_service.go -destination=mock_guest_list_service.go -package=guest_list; \


#	mockgen -source=guest_list_handler.go -destination=mock_guest_list_handler.go -package=guest_list; \


#mockgen -source=seats_service.go -destination=mock_seats_service.go -package=empty_seats
#mockgen -source=seats_handler.go -destination=mock_seats_handler.go -package=empty_seats

