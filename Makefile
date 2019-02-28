#Challenge Makefile

PROJECTNAME=$(shell basename "$(pwd)")
# Go related variables.
GOBASE=$(shell pwd)
$(shell export GOPATH=$(GOPATH):$(GOBASE)/vendor:$(GOBASE))
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will store the server process id when it's running on development mode
PID=/tmp/.$(PROJECTNAME)-api-server.pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## setup: Install dependencies to run project.
setup: go-clean go-get

## start: Start server.
start:
	docker-compose -f "docker-compose.yml" up -d --build
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run main.go

## check: Run tests.
check: go-test

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

go-test:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test -race -v ./...

cover:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test -race -v -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go tool cover -html=coverage.txt

help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
