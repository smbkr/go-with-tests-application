.PHONY: integration-test
integration-test: db
	go test -v ./...

.PHONY: db
db:
	docker-compose down
	docker-compose up -d
