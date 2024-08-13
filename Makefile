general:
	@echo === Available makefile commands:
	@echo run - run the service
	@echo install (alias - i) - install project dependencies
	@echo gen - generate documentation and mock files
	@echo upgrade - upgrade dependencies
	@echo clean - clean mod files
	@echo start - start docker compose with rebuild
	@echo up - start docker compose without rebuild
	@echo stop - stop docker container
	@echo test-unit - start tests excluding integration tests
	@echo test - start all tests with coverage profile
	@echo test-verbose - start tests with verbose flag

run: gen test-verbose start
	
i: install

install:
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@$(MAKE) clean

gen:
	@swag init --parseDependency -d ./internal/endpoint -g ../app/app.go -o ./docs
	@go generate ./...

upgrade: clean i
	@go get -u ./... && go mod tidy

clean:
	@go mod tidy

start:
	@docker compose up --build --timestamps --wait --wait-timeout 1800 --remove-orphans -d

stop:
	@docker compose stop

up:
	@docker compose up --timestamps --wait --wait-timeout 1800 --remove-orphans -d

test-unit: test-folder-creation gen
	@go test ./... -v -short

test-verbose: test-folder-creation gen
	@go test ./... -v

test: test-folder-creation gen
	@go test ./... -coverprofile=tests/coverage -coverpkg=./... && go tool cover -func=tests/coverage -o tests/coverage.func && go tool cover -html=tests/coverage -o tests/coverage.html

test-folder-creation:
ifeq ($(OS),Windows_NT)
	-@mkdir tests
else
	-@mkdir -p tests
endif