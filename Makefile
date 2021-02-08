.PHONY: build-app
	go mod tidy
	go mod download

.PHONY: start-local-db
	docker-compose -f