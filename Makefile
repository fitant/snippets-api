ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
APP=snippet-api
APP_VERSION:="1.0"
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"

deps:
	go mod download

compile:
	mkdir out
	go build -o $(APP_EXECUTABLE) -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" src/*.go

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

tidy:
	go mod tidy

serve: fmt vet
	env $(cat dev.env | xargs) go run cmd/*.go

build: deps compile
