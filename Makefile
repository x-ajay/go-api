POSTGRES_IMAGE=postgres:12
POSTGRES_CONTAINER=postgresdb
POSTGRES_USER=root
POSTGRES_PASSWORD=secret
POSTGRES_DB=bank
POSTGRES_PORT=5432

MIGRATION_SUFFIX=$(shell date +'%d-%m-%Y.%H.%M.%S')
MIGRATION_NAME=unnamed-$(MIGRATION_SUFFIX).sql

postgres:
	@echo "running postgres using docker"
	docker run --name $(POSTGRES_CONTAINER) -p $(POSTGRES_PORT):5432 -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_DB) -d $(POSTGRES_IMAGE)
	@echo "postgres running."

postgres-clean:
	@echo "stopping postgres using docker"
	docker stop $(POSTGRES_CONTAINER)
	@echo "postgres stopped."
	@echo "removing postgres using docker"
	docker rm $(POSTGRES_CONTAINER)
	@echo "postgres removed."

createdb:
	@echo "creating database using docker"
	docker exec -it $(POSTGRES_CONTAINER) createdb -U $(POSTGRES_USER) $(POSTGRES_DB)
	@echo "database created."

dropdb:
	@echo "dropping database using docker"
	docker exec -it $(POSTGRES_CONTAINER) dropdb -U $(POSTGRES_USER) $(POSTGRES_DB)
	@echo "database dropped."

migrationcreate:
	@echo "creating migration"
	migrate create -ext sql -dir db/migrations -seq $(MIGRATION_NAME)

migrationup:
	@echo "running migration up"
	# modify the command so that it will not generate down migration
	migrate -path db/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose up

migrationdown:
	@echo "running migration down"
	migrate -path db/migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" -verbose down

sqlc:
	@echo "running sqlc"
	sqlc generate

test:
	@echo "running test"
	go test -v -cover ./...

run:
	@echo "running application"
	go run -mod=vendor main.go

build:
	@echo "building application"
	go build -mod=vendor -o bin/main main.go

vendor:
	@echo "running go mod vendor"
	go mod vendor

PHONY: postgres postgres-clean createdb dropdb migrationup migrationdown sqlc test run build vendor