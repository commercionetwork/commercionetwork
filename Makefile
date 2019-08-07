PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
GOBIN ?= $(GOPATH)/bin
GOSUM := $(shell which gosum)

include Makefile.ledger

export GO111MODULE = on

all: build test

########################################
### Install

install: go.sum
	GO111MODULE=on go install -tags "$(build_tags)" ./cmd/cnd
	GO111MODULE=on go install -tags "$(build_tags)" ./cmd/cncli

########################################
### Build

build: go.sum
	GO111MODULE=on go build -o "build/cnd" -tags "$(build_tags)" ./cmd/cnd/main.go
	GO111MODULE=on go build -o "build/cncli" -tags "$(build_tags)" ./cmd/cncli/main.go

build-darwin: go.sum
	env GO111MODULE=on GOOS=darwin GOARCH=386 go build -o ./build/darwin/cncli-darwin-386 -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=darwin GOARCH=386 go build -o ./build/darwin/cnd-darwin-386 -tags "$(build_tags)" ./cmd/cnd/main.go
	env GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/cncli-darwin-amd64 -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/cnd-darwin-amd64 -tags "$(build_tags)" ./cmd/cnd/main.go

build-linux: go.sum
	env GO111MODULE=on GOOS=linux GOARCH=386 go build -o ./build/linux/cncli-linux-386 -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=linux GOARCH=386 go build -o ./build/linux/cnd-linux-386 -tags "$(build_tags)" ./cmd/cnd/main.go
	env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ./build/linux/cncli-linux-amd64 -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ./build/linux/cnd-linux-amd64 -tags "$(build_tags)" ./cmd/cnd/main.go

build-windows: go.sum
	env GO111MODULE=on GOOS=windows GOARCH=386 go build -o ./build/windows/cncli-windows-386.exe -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=windows GOARCH=386 go build -o ./build/windows/cnd-windows-386.exe -tags "$(build_tags)" ./cmd/cnd/main.go
	env GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ./build/windows/cncli-windows-amd64.exe -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ./build/windows/cnd-windows-amd64.exe -tags "$(build_tags)" ./cmd/cnd/main.go

build-all: go.sum
	make build-darwin
	make build-linux
	make build-windows

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

########################################
### Testing

test: test_unit

test_unit:
	@VERSION=$(VERSION) go test -mod=readonly $(PACKAGES_NOSIMULATION) -tags='ledger test_ledger_mock'