build:
	CGO_ENABLED=0 go build -o build/location-app .

run:
	build/location-app $(ARGS)

tests:
	go vet ./...
	go test -cover -race ./...

docs:
	swag init -g api/http/server.go

lint:
	golangci-lint run

migrate.up:
	migrate -database $(POSTGRESQL_URL) -path repositories/sql/migrations up

migrate-docker-compose.up:
	make migrate.up POSTGRESQL_URL=postgres://postgres:password@localhost:5432/location?sslmode=disable

migrate.down:
	migrate -database $(POSTGRESQL_URL) -path repositories/sql/migrations down

migrate-docker-compose.down:
	make migrate.down POSTGRESQL_URL=postgres://postgres:password@localhost:5432/location?sslmode=disable

jwt-token:
	jwt encode --alg HS256 --secret secret --sub $(shell uuidgen)

.PHONY: build run tests docs lint db.migrate.up db.migrate.down
