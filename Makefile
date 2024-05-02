DB_NAME=modak-test-db
DB_PORT=5432
MIGRATION_DIR=storage/migrations
ROUTE="host=localhost user=postgres password=password dbname=${DB_NAME} port=${DB_PORT} sslmode=disable"

DB_NAME_TEST=modak-test-db-test
DB_PORT_TEST=5469

define setup_env
    $(eval ENV_FILE := .env.test_integration)
    @echo " - setup env $(ENV_FILE)"
    $(eval include .env.test_integration)
    $(eval export)
endef

## get extra arguments and filter out commands from args
args = $(filter-out $@,$(MAKECMDGOALS))

.PHONY: all test unit_test e2e_test

test:

	echo "Starting test environment"
	$(call setup_env)
	make unit_test
	make e2e_test

unit_test:
	echo "/////////////////////////////////Starting unit test environment/////////////////////////////////"
	cd ./repository && go test ./... -coverprofile coverage.out -covermode count &&  go tool cover -func coverage.out | grep total | awk '{print $3}'
	cd ./usecase && go test ./... -coverprofile coverage.out -covermode count &&  go tool cover -func coverage.out | grep total | awk '{print $3}'
	cd ./repository && go tool cover -html=coverage.out -o cover.html
	cd ./usecase && go tool cover -html=coverage.out -o cover.html


e2e_test:
	echo "Starting test environment"
	$(call setup_env)
	echo "/////////////////////////////////Deleting fixtures/////////////////////////////////"
	rm -rf ./test/fixtures
	echo "/////////////////////////////////Starting E2E Test/////////////////////////////////"
	docker compose -f docker-compose-test.yaml up --build -d
	cd ./test &&  go test ./... || true
	docker compose -f docker-compose-test.yaml down
	echo "/////////////////////////////////Ending E2E Test/////////////////////////////////"


## default that allows accepting extra args
%:
    @:

.PHONY: migration
migration:
	goose -dir ${MIGRATION_DIR} create $(call args,defaultstring) sql
.PHONY: migration
migration-go:
	goose -dir ${MIGRATION_DIR} create $(call args,defaultstring) go

migrate-status:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} status

migrate-up:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} up
migrate-seeds:
	./seeds/goose-custom -dir seeds up -dbstring ${ROUTE}

migrate-down:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} down

migrate-rollback:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} reset

migrate-reset:
	goose -dir ${MIGRATION_DIR} postgres ${ROUTE} reset

mocks:
	mockery --dir=domain --output=domain/mocks --outpkg=mocks --all
